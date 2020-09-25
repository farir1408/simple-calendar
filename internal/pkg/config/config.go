package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"

	"gopkg.in/yaml.v2"
)

// EnvPrefix ...
const EnvPrefix = "CALENDAR_API"

// AppConfig ...
type AppConfig struct {
	Host             string `yaml:"CALENDAR_API_HOST",envconfig:"CALENDAR_API_HOST"`
	Port             int    `yaml:"CALENDAR_API_PORT",envconfig:"CALENDAR_API_PORT"`
	DebugPort        int    `yaml:"CALENDAR_API_DEBUG_PORT",envconfig:"CALENDAR_API_DEBUG_PORT"`
	DatabaseAddr     string `yaml:"CALENDAR_API_DATABASE_ADDR",envconfig:"CALENDAR_API_DATABASE_ADDR"`
	DatabaseDBName   string `yaml:"CALENDAR_API_DATABASE_DB_NAME",envconfig:"CALENDAR_API_DATABASE_DB_NAME"`
	DatabaseUser     string `yaml:"CALENDAR_API_DATABASE_USER",envconfig:"CALENDAR_API_DATABASE_USER"`
	DatabasePassword string `yaml:"CALENDAR_API_DATABASE_PASSWORD",envconfig:"CALENDAR_API_DATABASE_PASSWORD"`
	LogLvl           string `yaml:"CALENDAR_API_LOG_LEVEL",envconfig:"CALENDAR_API_LOG_LEVEL"`
}

// New ...
func New() *AppConfig {
	// If local-config-path didn't enable, get config from env.
	localConfigEnable := flag.Bool("run-local", false, "enable local config file")
	configPath := flag.String("config-path", "./configs/app/local.yaml", "path to config file")
	flag.Parse()
	if *localConfigEnable {
		return newConfigFromFile(*configPath)
	}

	return newConfigFromEnv()
}

func newConfigFromFile(filePath string) *AppConfig {
	c := &AppConfig{}
	data, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		fmt.Printf("can't get config from file: %+v\n", err)
		return nil
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil
	}
	return c
}

func newConfigFromEnv() *AppConfig {
	c := &AppConfig{}
	err := envconfig.Process(EnvPrefix, c)
	if err != nil {
		fmt.Printf("can't get config from env: %+v\n", err)
		return nil
	}
	return c
}
