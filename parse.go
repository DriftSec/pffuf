package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func getCommands(curFile string) {
	jsonFile, err := os.Open(curFile)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var output FfufOutput
	json.Unmarshal(byteValue, &output)

	commands = append(commands, output.CommandLine)

}

func getEndpoints(curFile string) {
	var navResults NavResults

	jsonFile, err := os.Open(curFile)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var output FfufOutput
	json.Unmarshal(byteValue, &output)

	for i := 0; i < len(output.Results); i++ {
		navResults.URL = output.Results[i].URL
		navResults.Endpoint = output.Results[i].Input.FUZZ
		navResults.Status = output.Results[i].Status
		navResults.Length = output.Results[i].Length
		navResults.Words = output.Results[i].Words
		navResults.Lines = output.Results[i].Lines
		results = append(results, navResults)
		origresults = append(origresults, navResults)
	}

}
