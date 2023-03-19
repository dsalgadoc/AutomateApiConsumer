package providers

import (
	"fmt"
	"myApiController/configs"
	"myApiController/domain"
	infrastructure "myApiController/infrastructure/client"
)

var clients = map[string]func(c configs.Client) (domain.DataRowClient, error){
	configs.EngineClientType: buildEngineHttpClient,
}

func GetDataRowClient(c configs.Client) (domain.DataRowClient, error) {
	client, exists := clients[c.Name]
	if !exists {
		return nil, fmt.Errorf("unable to build %s client", c.Name)
	}

	return client(c)
}

func buildEngineHttpClient(c configs.Client) (domain.DataRowClient, error) {
	return infrastructure.NewEngineClient(c.Path), nil
}
