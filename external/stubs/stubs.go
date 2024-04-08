package stubs

import (
	"assignment-2/external"
	"encoding/json"
	"fmt"
	"os"
)

func readStubFile(filePath string) []byte {
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	return file
}

func MeteoStub() external.MeteoObject {
	var meteoData external.MeteoObject
	meteoFile := readStubFile("./data/meteoWeather.json")
	err := json.Unmarshal(meteoFile, &meteoData)
	if err != nil {
		fmt.Printf("File parse error: %v\n", err)
		os.Exit(1)
	}
	return meteoData
}
