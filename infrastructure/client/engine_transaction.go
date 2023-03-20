package client

import (
	"encoding/json"
	"io"
	"myApiController/domain"
	"myApiController/domain/model"
	"net/http"
)

type engineClient struct {
	url        string
	HttpClient *http.Client
	headers    http.Header
}

func NewEngineClient(path string, headers map[string]string, httpClient http.Client) domain.DataRowClient {
	header := http.Header{}
	for key, value := range headers {
		header.Add(key, value)
	}

	return &engineClient{
		url:        path,
		HttpClient: &httpClient,
		headers:    header,
	}
}

func (e *engineClient) DoRequest(params map[string]string) (domain.DataExchange, error) {
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

	resp, err := e.HttpClient.Do(req)
	if err != nil {
		return model.EngineResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.EngineRequest{}, err
	}

	var engineResponse model.EngineResponse
	err = json.Unmarshal(respBody, &engineResponse)
	if err != nil {
		return engineResponse, err
	}

	return engineResponse, nil
}
