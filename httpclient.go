package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var (
	httpClient *http.Client
)

const (
	MaxIdleConns        int = 100
	MaxIdleConnsPerHost int = 100
	IdleConnTimeout     int = 90
)

func init() {
	httpClient = createHTTPClient()

}
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 20 * time.Second,
	}
	return client
}

func HttpRequest(url, method string, body io.Reader, cookies []*http.Cookie) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	//req.Header.Set("Content-Type", "application/json")
	for _, v := range cookies {
		fmt.Printf("cookie name: %v, cookie value: %v", v.Name, v.Value)
		req.AddCookie(v)
		//req.Header.Set(v.Name, v.Value)
	}
	response, err := httpClient.Do(req)
	if err != nil && response == nil {
		return nil, err

	}

	if response.Body != nil {
		defer response.Body.Close()
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request url: %v failed, ResponseBody: %v", url, string(responseBody))
	}
	return responseBody, nil

}
