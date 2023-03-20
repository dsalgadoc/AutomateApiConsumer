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

type engineClientTestScenario struct {
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
	s := startEngineClientTestScenario(t)
	headers := map[string]string{"hello": "world"}
	params := map[string]string{"lorem": "ipsum"}
	s.givenAServerOk(headers)
	s.andSomeParams(params)
	s.whenDoingRequest()
	s.thenThereIsNoError()
}

func TestRequestWithBadResponse(t *testing.T) {
	s := startEngineClientTestScenario(t)
	headers := map[string]string{"hello": "world"}
	params := map[string]string{"lorem": "ipsum"}
	s.givenAServerBadResponse(headers)
	s.andSomeParams(params)
	s.andAnExpectedError(fmt.Errorf("invalid character 'H' looking for beginning of value"))
	s.whenDoingRequest()
	s.thenThereIsAnError()
}

func TestErrorProcessingRequest(t *testing.T) {
	s := startEngineClientTestScenario(t)
	headers := map[string]string{"hello": "world"}
	params := map[string]string{"lorem": "ipsum"}
	s.givenAServerWithError(headers)
	s.andSomeParams(params)
	s.andAnExpectedError(fmt.Errorf(`Get "http://test.com/GET?lorem=ipsum": something went wrong`))
	s.whenDoingRequest()
	s.thenThereIsAnError()
}

func TestErrorEOFResponse(t *testing.T) {
	s := startEngineClientTestScenario(t)
	headers := map[string]string{"hello": "world"}
	params := map[string]string{"lorem": "ipsum"}
	s.givenAServerWithErrorResponseEOF(headers)
	s.andSomeParams(params)
	s.andAnExpectedError(fmt.Errorf("invalid character 'H' looking for beginning of value"))
	s.whenDoingRequest()
	s.thenThereIsAnError()
}

/*-- steps --*/
func startEngineClientTestScenario(t *testing.T) *engineClientTestScenario {
	t.Parallel()
	return &engineClientTestScenario{
		test: t,
	}
}

func (e *engineClientTestScenario) givenAServerOk(headers map[string]string) {
	e.mockClient = http.Client{Transport: &mockClientOk{}}
	e.client = NewEngineClient("http://test.com/GET", headers, e.mockClient)
}

func (e *engineClientTestScenario) givenAServerBadResponse(headers map[string]string) {
	e.mockClient = http.Client{Transport: &mockClientBadResponse{}}
	e.client = NewEngineClient("http://test.com/GET", headers, e.mockClient)
}

func (e *engineClientTestScenario) givenAServerWithError(headers map[string]string) {
	e.mockClient = http.Client{Transport: &mockClientError{}}
	e.client = NewEngineClient("http://test.com/GET", headers, e.mockClient)
}

func (e *engineClientTestScenario) givenAServerWithErrorResponseEOF(headers map[string]string) {
	e.mockClient = http.Client{Transport: &mockClientResponseEOF{}}
	e.client = NewEngineClient("http://test.com/GET", headers, e.mockClient)
}

func (e *engineClientTestScenario) andSomeParams(params map[string]string) {
	e.params = params
}

func (e *engineClientTestScenario) andAnExpectedError(err error) {
	e.expectedErr = err
}

func (e *engineClientTestScenario) whenDoingRequest() {
	e.result, e.err = e.client.DoRequest(e.params)
}

func (e *engineClientTestScenario) thenThereIsNoError() {
	assert.NoError(e.test, e.err)
}

func (e *engineClientTestScenario) thenThereIsAnError() {
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
