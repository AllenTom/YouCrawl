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
	Process(item interface{}, store GlobalStore) error
}

type ImageDownloadItem struct {
	Urls []string
}
type ImageDownloadPipeline struct {
	// get store folder
	//
	//./download/image by default
	GetStoreFileFolder func(item interface{}, store GlobalStore) string
	// get save filename
	//
	// same name with image,by default
	GetSaveFileName func(item interface{}, store GlobalStore, rawURL string) string
	// get urls
	//
	//if the type of Item is ImageDownloadItem, no need to specify
	GetUrls func(item interface{}, store GlobalStore) []string
	// maximum number of concurrent downloads
	MaxDownload int
	// request middlewares to use
	Middlewares []Middleware
	// call on each image downloaded complete
	OnImageDownloadComplete func(item interface{}, store GlobalStore, url string, downloadFilePath string)
	// call on all image download, regardless of whether all image download is successful
	OnDone func(item interface{}, store GlobalStore)
}

func (i *ImageDownloadPipeline) Process(item interface{}, store GlobalStore) error {
	logger := logrus.WithField("scope", "Image Download Pipeline")
	// prepare
	rawDownloadURLs, ok := item.(ImageDownloadItem)

	// not download image item,use pipeline get url method
	var urls []string
	if !ok {
		if i.GetUrls != nil {
			urls = i.GetUrls(item, store)
		}
	} else {
		urls = rawDownloadURLs.Urls
	}

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
				if i.OnImageDownloadComplete != nil {
					i.OnImageDownloadComplete(item, store, downloadURL, saveFilePath)
				}
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
	if i.OnDone != nil {
		i.OnDone(item, store)
	}
	return nil
}

type ChannelPipelineToken interface {
	GetToken() string
}
type ChannelPipeline struct {
	ChannelMapping sync.Map
}

func (p *ChannelPipeline) Process(item interface{}, _ GlobalStore) error {
	if channelToken,isToken := item.(ChannelPipelineToken);isToken {
		tokenString := channelToken.GetToken()

		if rawChannel,isExist := p.ChannelMapping.Load(tokenString);isExist {
			channel := rawChannel.(chan interface{})
			channel <- item
		}
	}
	return nil
}