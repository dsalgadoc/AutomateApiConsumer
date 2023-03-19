package application

import (
	"fmt"
	"myApiController/configs"
	"myApiController/domain"
)

type DataProcessor struct {
	config    configs.Config
	inputter  domain.DataInputter
	outputter domain.DataOutputter
	rowClient domain.DataRowClient
}

func NewDataProcessor(c configs.Config,
	i domain.DataInputter,
	o domain.DataOutputter,
	r domain.DataRowClient) DataProcessor {
	return DataProcessor{
		config:    c,
		inputter:  i,
		outputter: o,
		rowClient: r,
	}
}

func (dp *DataProcessor) Do() {
	filepath := dp.config.IO.FolderLocation
	data, err := dp.inputter.Invoke(filepath + dp.config.IO.InputFileName + dp.inputter.InputterExtension())
	if err != nil {
		panic(fmt.Errorf("data inputter error: %w", err))
	}

	dataReturned := dp.getDataFromRegisteredClient(data)

	err = dp.outputter.Write(filepath+dp.outputter.OutputterFilename(), dataReturned)
	if err != nil {
		panic(fmt.Errorf("data outputter error: %w", err))
	}
}

func (dp *DataProcessor) getDataFromRegisteredClient(data domain.Table) []domain.DataExchange {
	dataReturned := []domain.DataExchange{}
	for _, row := range data.Rows {
		params := dp.rowToParams(data.Headers, row)
		rowProcessed, err := dp.rowClient.DoRequest(params)
		if err != nil {
			continue
		}
		dataReturned = append(dataReturned, rowProcessed)
	}
	return dataReturned
}

func (dp *DataProcessor) rowToParams(headers []string, row []string) map[string]string {
	params := make(map[string]string)
	for i, cell := range row {
		params[headers[i]] = cell
	}
	return params
}
