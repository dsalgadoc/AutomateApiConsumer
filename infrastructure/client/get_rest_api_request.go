package client

import (
	"encoding/json"
	"io"
	"myApiController/domain"
	"myApiController/domain/model"
	"net/http"
)

type getRestApi struct {
	url        string
	HttpClient *http.Client
	headers    http.Header
}

func NewGetRestApi(path string, headers map[string]string, httpClient http.Client) domain.DataRowClient {
	header := http.Header{}
	for key, value := range headers {
		header.Add(key, value)
	}

	return &getRestApi{
		url:        path,
		HttpClient: &httpClient,
		headers:    header,
	}
}

func (e *getRestApi) DoRequest(params map[string]string) (domain.DataExchange, error) {
	req, err := http.NewRequest(http.MethodGet, e.url, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()
	req.Header = e.headers

	var response model.ApiResponse
	resp, err := e.HttpClient.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
