package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var myClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func DoPost(url string, sessionId string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("responsecompressed", "0")
	req.Header.Set("requestcompressed", "0")
	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", sessionId))
	return myClient.Do(req)
}

func DoGet(url string, sessionId string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("responsecompressed", "0")
	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", sessionId))
	return myClient.Do(req)
}

func DoGetCompressed(url string, sessionId string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("responsecompressed", "1")
	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", sessionId))
	return myClient.Do(req)
}
