package main

import (
	"fmt"
	"myApiController/cmd/application"
	"myApiController/configs"
	"os"
)

// registered IO and client
var (
	IOsType    = []string{configs.CsvIoType, configs.JsonIoType}
	ClientType = []string{configs.EngineClientType}
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
	if !application.CheckArgumentOnSlice(client, ClientType) {
		panic(fmt.Errorf("invalid client, the valid types are %+v", ClientType))
	}

	var app = application.BuildApplication(inputter, outputter, client)
	app.DataProcessor.Do()
}
