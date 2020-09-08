package youcrawl

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func RequestWithURL(requestURL string) error {
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	fmt.Println(fmt.Sprintf("%s [%d]", requestURL, resp.StatusCode))

	return nil
}
