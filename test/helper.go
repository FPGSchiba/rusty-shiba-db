package test

import (
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"net/http/httptest"
	router2 "rsdb/src/router"
)

func request(method string, url string, body io.Reader) (*http.Response, error) {
	router := router2.GetRouter()
	w := httptest.NewRecorder()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	router.ServeHTTP(w, req)
	return w.Result(), nil
}

func getResponseBody(resp *http.Response) (map[string]interface{}, error) {
	var responseBody map[string]interface{}
	var bodyContent []byte
	_, err := resp.Body.Read(bodyContent)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyContent, &responseBody)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}
