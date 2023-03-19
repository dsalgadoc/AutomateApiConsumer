package application

import (
	"fmt"
	"myApiController/application"
	"myApiController/cmd/builders"
	"myApiController/cmd/providers"
	"myApiController/configs"
	"myApiController/domain"
)

type (
	Application struct {
		DataProcessor application.DataProcessor
	}
)

func BuildApplication(inputterType, outputterType, clientType string) *Application {
	appConfig := getConfiguration()
	inputter, err := buildInputter(inputterType)
	if err != nil {
		panic(fmt.Errorf("error building inputter: %w", err))
	}

	outputter, err := buildOutputter(outputterType)
	if err != nil {
		panic(fmt.Errorf("error building outputter: %w", err))
	}

	client, err := buildClients(clientType, appConfig)
	if err != nil {
		panic(fmt.Errorf("error building clients: %w", err))
	}

	return &Application{
		DataProcessor: builders.BuildDataProcessor(appConfig, inputter, outputter, client),
	}
}

func getConfiguration() configs.Config {
	appConfig, err := configs.LoadConfig("./api/configs/config.yaml")
	if err != nil {
		panic(fmt.Errorf("error getting configuration: %w", err))
	}
	return appConfig
}

func buildInputter(iType string) (domain.DataInputter, error) {
	return providers.GetDataInputter(iType)
}

func buildOutputter(oType string) (domain.DataOutputter, error) {
	return providers.GetDataOutputter(oType)
}

func buildClients(cType string, c configs.Config) (domain.DataRowClient, error) {
	var client configs.Client
	for _, cli := range c.Clients {
		if cType == cli.Name {
			client = cli
		}
	}
	return providers.GetDataRowClient(client)
}
