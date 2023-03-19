package builders

import (
	"myApiController/application"
	"myApiController/configs"
	"myApiController/domain"
)

func BuildDataProcessor(
	config configs.Config,
	inputter domain.DataInputter,
	outputter domain.DataOutputter,
	client domain.DataRowClient,
) application.DataProcessor {
	return application.NewDataProcessor(config, inputter, outputter, client)
}
