package mocks

import (
	"github.com/stretchr/testify/mock"
	"myApiController/domain"
)

type DataRowClientMock struct {
	mock.Mock
}

func (m *DataRowClientMock) DoRequest(params map[string]string) (domain.DataExchange, error) {
	args := m.Called(params)
	return args.Get(0).(domain.DataExchange), args.Error(1)
}
