package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const ErrMsg = "Couldn't load configuration, cannot start. Terminating. Error:"

// loadConfiguration for example http://configserver:8888/accountservice/test
func loadConfiguration(configServerURL string, appName string, profile string) {
	url := fmt.Sprintf("%s/%s/%s", configServerURL, appName, profile)
	log.Printf("Loading config from %s\n", url)
	body, err := fetchConfiguration(url)
	if err != nil {
		log.Fatalf("%s %s", ErrMsg, err.Error())
	}
	parseConfiguration(body)
}

// Make HTTP request to fetch configuration from config server
func fetchConfiguration(url string) ([]byte, error) {
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("%s %s", ErrMsg, err.Error())
	}
	req.SetBasicAuth(viper.GetString("CLOUD_USERNAME"), viper.GetString("CLOUD_PASSWORD"))
	resp, errDo := client.Do(req)
	if errDo != nil {
		log.Fatalf("%s %s", ErrMsg, err.Error())
	}
	body, errDo := ioutil.ReadAll(resp.Body)
	return body, errDo
}

// Pass JSON bytes into struct and then into Viper
func parseConfiguration(body []byte) {
	var cloudConfig springCloudConfig
	err := json.Unmarshal(body, &cloudConfig)
	if err != nil {
		log.Fatalf("Cannot parse configuration, message: %s", err.Error())
	}
	for key, value := range cloudConfig.PropertySources[0].Source {
		viper.Set(key, value)
		log.Printf("Loading config property %v => %v\n", key, value)
	}
	if viper.IsSet("server_name") {
		log.Printf("Successfully loaded configuration for service %s\n", viper.GetString("server_name"))
	}
}

// Structs having same structure as response from Spring Cloud Config
type springCloudConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	PropertySources []propertySource `json:"propertySources"`
}

type propertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}
