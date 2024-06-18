package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func HttpJsonGet(requestUrl string) (buf []byte, err error) {
	var client = http.Client{}
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return
	}
	//req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("do http request error code: %d", resp.StatusCode)
	}
	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

func HttpJsonPost(requestUrl string, body []byte) (buf []byte, err error) {
	var client = http.Client{}
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("do http request error code: %d", resp.StatusCode)
	}
	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}
