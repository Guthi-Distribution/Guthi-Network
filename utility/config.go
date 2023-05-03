package utility

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Name                       string   `json:"name"`
	Address                    string   `json:"address"`
	Port                       int      `json:"port"`
	InitialConnectionAddresses []string `json:"connect_to"`
	ShareNodes                 bool     `json:"share_nodes"`
}

func LoadConfiguration(filename string) Config {
	config := Config{}
	config.ShareNodes = false

	configFile, err := os.Open(filename)
	defer configFile.Close()
	if err != nil {
		config.Name = ""
		config.Address = ""
		log.Fatalf("Config file opening error: %s\n", err)
		return config
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)

	if err != nil {
		log.Fatalf("Invalid Json for configuration: %s\n", err)
	}

	return config
}
