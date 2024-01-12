package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// General Function to load any config file as interface type
func LoadConfig(path string) []interface{} {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Read the contents of the file
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	// Create a variable to hold the unmarshalled data
	var config []interface{}

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	return config
}

// Load DB config file
func LoadDBConfig(path string, dType any) (interface{}, error) {
	// DB struct
	var config interface{}

	// open file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return config, err
	}
	defer file.Close()

	// decode file contents into the data struct
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func LoadCloudConfig() {

}
