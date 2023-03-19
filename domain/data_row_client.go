package domain

type DataRowClient interface {
	DoRequest(params map[string]string) (DataExchange, error)
}
