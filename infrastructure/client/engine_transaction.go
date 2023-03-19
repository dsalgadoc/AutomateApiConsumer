package client

import (
	"encoding/json"
	"myApiController/domain"
	"myApiController/domain/model"
	"net/http"
	"time"
)

type engineClient struct {
	url        string
	HttpClient *http.Client
}

func NewEngineClient(path string) domain.DataRowClient {
	return &engineClient{
		url: path,
		HttpClient: &http.Client{
			Timeout: time.Second * 1000,
		},
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

	resp, err := e.HttpClient.Do(req)
	if err != nil {
		return model.EngineResponse{}, err
	}
	defer resp.Body.Close()

	var engineResponse model.EngineResponse
	err = json.NewDecoder(resp.Body).Decode(&engineResponse)
	if err != nil {
		return engineResponse, err
	}

	return engineResponse, nil
}
