package youcrawl

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Pipeline interface {
	Process(item *Item, store GlobalStore) error
}

type ImageDownloadPipeline struct {
	GetStoreFileFolder func(item *Item, store GlobalStore) string
	GetSaveFileName    func(item *Item, store GlobalStore, rawURL string) string
	MaxDownload        int
	Middlewares        []Middleware
}

func (i *ImageDownloadPipeline) Process(item *Item, store GlobalStore) error {
	logger := logrus.WithField("scope", "Image Download Pipeline")
	// prepare
	rawDownloadURLs, err := item.GetValue("downloadImgURLs")
	if err != nil {
		return err
	}
	urls := rawDownloadURLs.([]string)

	maxDownload := i.MaxDownload

	// use default max download
	if i.MaxDownload == 0 {
		maxDownload = 3
	}

	workChan := make(chan struct{}, maxDownload)
	for idx := 0; idx < maxDownload; idx++ {
		workChan <- struct{}{}
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))
	for idx, rawDownloadUrl := range urls {
		<-workChan
		storePath := filepath.Join("./", "download", "images")
		if i.GetStoreFileFolder != nil {
			storePath = i.GetStoreFileFolder(item, store)
		}

		saveFileName := ""
		if i.GetSaveFileName != nil {
			saveFileName = i.GetSaveFileName(item, store, rawDownloadUrl)
		}
		if saveFileName == "" {
			downloadURL, err := url.Parse(rawDownloadUrl)
			if err != nil {
				return err
			}
			urlParts := strings.Split(downloadURL.Path, "/")
			saveFileName = urlParts[len(urlParts)-1]
		}

		savePath := filepath.Join(storePath, saveFileName)

		go func(workChan chan struct{}, wg *sync.WaitGroup, saveFilePath string, downloadURL string, idx int) {
			defer func() {
				workChan <- struct{}{}
				wg.Done()
			}()
			err := os.MkdirAll(filepath.Dir(saveFilePath), os.ModePerm)
			if err != nil {
				logger.Error(err)
				return
			}

			req, err := http.NewRequest("GET", downloadURL, nil)
			if err != nil {
				logger.Error(err)
				return
			}
			client := &http.Client{}
			if i.Middlewares != nil {
				for _, middleware := range i.Middlewares {
					middleware.Process(client, req, &Context{})
				}
			}

			resp, err := client.Do(req)
			if err != nil {
				logger.Error(err)
				return
			}

			defer resp.Body.Close()

			file, err := os.Create(saveFilePath)
			if err != nil {
				logger.Error(err)
				return
			}
			defer file.Close()

			_, err = io.Copy(file, resp.Body)
			if err != nil {
				logger.Error(err)
				return
			}

		}(workChan, &wg, savePath, rawDownloadUrl, idx)
	}

	wg.Wait()
	return nil
}
