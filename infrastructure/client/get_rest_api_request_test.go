package client

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"myApiController/domain"
	"net/http"
	"strings"
	"testing"
)

const (
	EXAMPLE_SIMPLE_URL  = "http://test.com/GET"
	EXAMPLE_WITH_PARAMS = "http://test.com/GET/{param1}/{param2}"
)

type getRestApiTestScenario struct {
	test        *testing.T
	client      domain.DataRowClient
	mockClient  http.Client
	params      map[string]string
	result      domain.DataExchange
	err         error
	expectedErr error
}

/*-- Test --*/
func TestRequestOk(t *testing.T) {
	s := startGetRestApiTestScenario(t)
	headers := map[string]string{"hello": "world"}
	params := map[string]string{"lorem": "ipsum"}
	s.givenAServerOk(EXAMPLE_SIMPLE_URL, headers)
	s.andSomeParams(params)
	s.whenDoingRequest()
	s.thenThereIsNoError()
}

func TestRequestOkWithParams(t *testing.T) {
	s := startGetRestApiTestScenario(t)
	headers := map[string]string{"hello": "world"}
	params := map[string]string{"lorem": "ipsum", "param1": "12345", "param2": "sample"}
	s.givenAServerOk(EXAMPLE_WITH_PARAMS, headers)
	s.andSomeParams(params)
	s.whenDoingRequest()
	s.thenThereIsNoError()
}

func TestRequestOkWithParamsNoProvided(t *testing.T) {
	s := startGetRestApiTestScenario(t)
	headers := map[string]string{"hello": "world"}
	params := map[string]string{}
	s.givenAServerOk(EXAMPLE_WITH_PARAMS, headers)
	s.andSomeParams(params)
	s.whenDoingRequest()
	s.thenThereIsNoError()
}

func TestRequestWithBadResponse(t *testing.T) {
	s := startGetRestApiTestScenario(t)
	headers := map[string]string{"hello": "world"}
	params := map[string]string{"lorem": "ipsum"}
	s.givenAServerBadResponse(EXAMPLE_SIMPLE_URL, headers)
	s.andSomeParams(params)
	s.andAnExpectedError(fmt.Errorf("invalid character 'H' looking for beginning of value"))
	s.whenDoingRequest()
	s.thenThereIsAnError()
}

func TestErrorProcessingRequest(t *testing.T) {
	s := startGetRestApiTestScenario(t)
	headers := map[string]string{"hello": "world"}
	params := map[string]string{"lorem": "ipsum"}
	s.givenAServerWithError(EXAMPLE_SIMPLE_URL, headers)
	s.andSomeParams(params)
	s.andAnExpectedError(fmt.Errorf(`Get "http://test.com/GET?lorem=ipsum": something went wrong`))
	s.whenDoingRequest()
	s.thenThereIsAnError()
}

func TestErrorEOFResponse(t *testing.T) {
	s := startGetRestApiTestScenario(t)
	headers := map[string]string{"hello": "world"}
	params := map[string]string{"lorem": "ipsum"}
	s.givenAServerWithErrorResponseEOF(EXAMPLE_SIMPLE_URL, headers)
	s.andSomeParams(params)
	s.andAnExpectedError(fmt.Errorf("invalid character 'H' looking for beginning of value"))
	s.whenDoingRequest()
	s.thenThereIsAnError()
}

/*-- steps --*/
func startGetRestApiTestScenario(t *testing.T) *getRestApiTestScenario {
	t.Parallel()
	return &getRestApiTestScenario{
		test: t,
	}
}

func (e *getRestApiTestScenario) givenAServerOk(path string, headers map[string]string) {
	e.mockClient = http.Client{Transport: &mockClientOk{}}
	e.client = NewGetRestApi(path, headers, e.mockClient)
}

func (e *getRestApiTestScenario) givenAServerBadResponse(path string, headers map[string]string) {
	e.mockClient = http.Client{Transport: &mockClientBadResponse{}}
	e.client = NewGetRestApi(path, headers, e.mockClient)
}

func (e *getRestApiTestScenario) givenAServerWithError(path string, headers map[string]string) {
	e.mockClient = http.Client{Transport: &mockClientError{}}
	e.client = NewGetRestApi(path, headers, e.mockClient)
}

func (e *getRestApiTestScenario) givenAServerWithErrorResponseEOF(path string, headers map[string]string) {
	e.mockClient = http.Client{Transport: &mockClientResponseEOF{}}
	e.client = NewGetRestApi(path, headers, e.mockClient)
}

func (e *getRestApiTestScenario) andSomeParams(params map[string]string) {
	e.params = params
}

func (e *getRestApiTestScenario) andAnExpectedError(err error) {
	e.expectedErr = err
}

func (e *getRestApiTestScenario) whenDoingRequest() {
	e.result, e.err = e.client.DoRequest(e.params)
}

func (e *getRestApiTestScenario) thenThereIsNoError() {
	assert.NoError(e.test, e.err)
}

func (e *getRestApiTestScenario) thenThereIsAnError() {
	assert.Error(e.test, e.err)
	assert.Equal(e.test, e.expectedErr.Error(), e.err.Error())
}

/*-- mocks --*/
type mockClientOk struct{}

func (t *mockClientOk) RoundTrip(req *http.Request) (*http.Response, error) {
	res := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("{}")),
	}
	return res, nil
}

type mockClientBadResponse struct{}

func (t *mockClientBadResponse) RoundTrip(req *http.Request) (*http.Response, error) {
	res := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("Hello, world")),
	}
	return res, nil
}

type mockClientError struct{}

func (t *mockClientError) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("something went wrong")
}

type mockClientResponseEOF struct{}

func (t *mockClientResponseEOF) RoundTrip(req *http.Request) (*http.Response, error) {
	body := "Hello, world"
	limit := int64(len(body)) - 1
	reader := io.LimitReader(strings.NewReader(body), limit)
	res := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(reader),
	}
	return res, nil
}
