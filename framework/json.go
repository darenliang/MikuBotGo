package framework

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// GetJsonString returns a json string from a url
func GetJsonString(rawURL string) string {
	var httpClient = &http.Client{Timeout: 10 * time.Second}
	escapedURL := url.QueryEscape(rawURL)
	finalURL, err := url.Parse(escapedURL)
	if err != nil {
		panic(err)
	}
	response, err := httpClient.Get(finalURL.Path)
	if err != nil {
		panic(err)
	}
	str, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	err = response.Body.Close()
	if err != nil {
		panic(err)
	}
	return string(str)
}
