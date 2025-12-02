package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Instances              []Instance `json:"instances"`
	ReconcileRetryInterval int        `json:"reconcileRetryInterval"`
	ListenAddress          string     `json:"listenAddress"`
	StoragePath            string     `json:"storagePath"`
}

type Instance struct {
	HomeServerName string `json:"homeServerName"`
	ApiAccessToken string `json:"apiAccessToken"`

	CorporalApiUrl      string `json:"corporalApiUrl"`
	CorporalAccessToken string `json:"corporalAccessToken"`
}

func Load(configPath string) (Config, error) {
	configString, ok := os.LookupEnv("MCMTP_CONFIG")
	if ok {
		return Parse([]byte(configString))
	}

	if configPath == "" {
		configPath = os.Getenv("MCMTP_CONFIG_PATH")
	}

	if configPath == "" {
		configPath = "config.json"
	}

	file, err := os.Open(configPath)
	if err != nil {
		logrus.Fatal(err)
	}

	configBytes, err := io.ReadAll(file)
	if err != nil {
		logrus.Fatal(err)
	}

	return Parse(configBytes)
}

func Parse(configBytes []byte) (Config, error) {
	var config Config

	if err := json.Unmarshal(configBytes, &config); err != nil {
		return Config{}, err
	}

	// set default. todo move somewhere?
	if config.ReconcileRetryInterval == 0 {
		config.ReconcileRetryInterval = 30
	}

	if err := Check(config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func Check(config Config) error {
	if config.ListenAddress == "" {
		return fmt.Errorf(`Config invalid: listenAddress must be set.`)
	}

	if config.StoragePath == "" {
		return fmt.Errorf(`Config invalid: storagePath must be set.`)
	}

	if len(config.Instances) == 0 {
		return fmt.Errorf(`Config invalid: At lease one instance (instances) needs to be configured.`)
	}

	for i, instance := range config.Instances {
		if instance.HomeServerName == "" {
			return fmt.Errorf(`Config invalid: Instance at index %d is missing a homeServerName.`, i)
		}

		if instance.ApiAccessToken == "" {
			return fmt.Errorf(`Config invalid: apiAccessToken must be set. Instance: %s`, instance.HomeServerName)
		}

		if instance.CorporalApiUrl == "" {
			return fmt.Errorf(`Config invalid: corporalApiUrl must be set. Instance: %s`, instance.HomeServerName)
		}

		if instance.CorporalAccessToken == "" {
			return fmt.Errorf(`Config invalid: corporalAccessToken must be set. Instance: %s`, instance.HomeServerName)
		}
	}

	return nil
}
