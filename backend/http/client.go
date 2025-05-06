package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var instance *Client

// Client struct to encapsulate HTTP client and methods
type Client struct {
	httpClient *http.Client
}

func NewClient(timeoutInSeconds uint16) {
	instance = &Client{
		&http.Client{
			Timeout: time.Duration(timeoutInSeconds) * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func getInstance() *Client {
	if instance == nil {
		log.Println("HTTP client not initialized, creating new client with default timeout")
		NewClient(10)
	}
	return instance
}

func DoPost(url string, sessionId string, body interface{}) (*http.Response, error) {
	myClient := getInstance().httpClient
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
	myClient := getInstance().httpClient
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("responsecompressed", "0")
	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", sessionId))
	return myClient.Do(req)
}

func DoGetCompressed(url string, sessionId string) (*http.Response, error) {
	myClient := getInstance().httpClient
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("responsecompressed", "1")
	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", sessionId))
	return myClient.Do(req)
}
