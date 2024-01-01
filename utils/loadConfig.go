package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// load the server configuration
func LoadSrvConfig() []interface{} {
	file, err := os.Open("./config/servers.json")
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
	var myData []interface{}

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(data, &myData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	return myData
}

func LoadCloudConfig() {

}