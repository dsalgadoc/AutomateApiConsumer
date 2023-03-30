package client

import (
	"encoding/json"
	"io"
	"myApiController/domain"
	"myApiController/domain/model"
	"net/http"
	"regexp"
)

const INNER_VARAIBLE = `\{(\w+)\}`

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
	path := e.generatePath(params)
	req, err := http.NewRequest(http.MethodGet, path, nil)
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

func (e *getRestApi) generatePath(params map[string]string) string {
	re := regexp.MustCompile(INNER_VARAIBLE)
	if !re.MatchString(e.url) {
		return e.url
	}
	return re.ReplaceAllStringFunc(e.url, func(match string) string {
		key := match[1 : len(match)-1]
		value, ok := params[key]
		if ok {
			return value
		}
		return match
	})
}
