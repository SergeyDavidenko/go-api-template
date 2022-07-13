package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

const ErrMsg = "Couldn't load configuration, cannot start. Terminating. Error:"

// loadConfiguration for example http://configserver:8888/accountservice/test
func loadConfiguration(configServerURL string, appName string, profile string) {
	if profile == "" {
		profile = "default"
	}
	if configServerURL == "" {
		log.Fatal("Config server url not set")
	}
	if appName == "" {
		log.Fatal("appName for config server not set")
	}
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
	client := http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("%s %s", ErrMsg, err.Error())
	}
	if viper.GetString("CLOUD_USERNAME") != "" && viper.GetString("CLOUD_PASSWORD") != "" {
		req.SetBasicAuth(viper.GetString("CLOUD_USERNAME"), viper.GetString("CLOUD_PASSWORD"))
	} else {
		log.Infof("skip auth config-server")
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("%s %s", ErrMsg, err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal("check cloud config url. Status code not eq 200. Status code is:", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// Pass JSON bytes into struct and then into Viper
func parseConfiguration(body []byte) {
	var cloudConfig springCloudConfig
	err := json.Unmarshal(body, &cloudConfig)
	if err != nil {
		log.Fatalf("cannot parse configuration, message: %s", err.Error())
	}
	if len(cloudConfig.PropertySources) > 0 {
		for key, value := range cloudConfig.PropertySources[0].Source {
			viper.Set(key, value)
			// log.Printf("Loading config property %v => %v\n", key, value)
		}
		if viper.IsSet("server_name") {
			log.Infof("successfully loaded configuration for service %s", viper.GetString("server_name"))
		}
	} else {
		log.Fatalf("cannot get config from cloud config service. Check if exist config.")
	}
}

// Pass JSON bytes into struct and them
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

func CloudConfigClient(appName string) {
	viper.AutomaticEnv()
	if viper.GetBool("CLOUD_CONFIG") {
		if viper.GetString("CLOUD_URL") == "" {
			log.Error("env CLOUD_URI not set")
		}
		if appName == "" {
			appName = viper.GetString("SERVICE_NAME")
		}
		loadConfiguration(viper.GetString("CLOUD_URL"), appName, "")
	} else {
		log.Info("loading config from file")
	}
}
