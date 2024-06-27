package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func FormatUrl(apiServer, path string) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(apiServer, "/"), strings.TrimLeft(path, "/"))
}

func HttpJsonGet(requestUrl string, headers map[string]string) (buf []byte, err error) {
	var client = http.Client{}
	log.Println("[http.get.url]", requestUrl)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return
	}
	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("X_CLIENT_ID", ClientId())
	if headers != nil {
		for h, v := range headers {
			req.Header.Set(h, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("do http request error code: %d", resp.StatusCode)
	}
	buf, err = io.ReadAll(resp.Body)
	log.Println("[http.post.resp]", string(buf))
	if err != nil {
		return
	}
	return
}

func HttpJsonPost(requestUrl string, body []byte, headers map[string]string) (buf []byte, err error) {
	var client = http.Client{}
	log.Println("[http.post.url]", requestUrl)
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("X_CLIENT_ID", ClientId())
	if headers != nil {
		for h, v := range headers {
			req.Header.Set(h, v)
		}
	}
	req.Header.Set("Authorization", ClientId())
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("do http request error code: %d", resp.StatusCode)
	}
	buf, err = io.ReadAll(resp.Body)
	log.Println("[http.post.resp]", string(buf))
	if err != nil {
		return
	}
	return
}

func HttpJsonGetUnmarshal[T any](requestUrl string, headers map[string]string, out *T) (resp *T, err error) {
	buf, err := HttpJsonGet(requestUrl, headers)
	if err != nil {
		return out, err
	}
	if err = json.Unmarshal(buf, out); err != nil {
		return out, err
	}
	return out, nil
}

func HttpJsonPostUnmarshal[T any](requestUrl string, body []byte, headers map[string]string, out *T) (resp *T, err error) {
	buf, err := HttpJsonPost(requestUrl, body, headers)
	if err != nil {
		return out, err
	}
	if err = json.Unmarshal(buf, out); err != nil {
		return out, err
	}
	return out, nil
}
