package main

import (
	"fmt"
	"myApiController/cmd/application"
	"myApiController/configs"
	"os"
)

// registered
var (
	IOsType = []string{configs.CsvIoType, configs.JsonIoType}
)

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		panic("no enough params")
	}

	inputter := args[0]
	if !application.CheckArgumentOnSlice(inputter, IOsType) {
		panic(fmt.Errorf("invalid inputter, the valid types are %+v", IOsType))
	}

	outputter := args[1]
	if !application.CheckArgumentOnSlice(outputter, IOsType) {
		panic(fmt.Errorf("invalid outputter, the valid types are %+v", IOsType))
	}

	client := args[2]

	var app = application.BuildApplication(inputter, outputter, client)

	registeredClientsNames := app.AppConfig.GetRegisteredClientsNames()
	if !application.CheckArgumentOnSlice(client, registeredClientsNames) {
		panic(fmt.Errorf("invalid client, the valid types are %+v", registeredClientsNames))
	}

	fmt.Printf("...Running request from (%s), recovery data from (%s) and writing to (%s)\n",
		client, inputter, outputter)
	app.DataProcessor.Do()
}
