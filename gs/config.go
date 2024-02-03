package gs

import (
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

type IConfiguration interface {
	GetActiveEnv() string
	GetGinRelease() bool
	GetGinAddr() string
	SolveDefaultValue()
	GetSnowFlakeConfig() SnowFlakeConfig
}

type SnowFlakeConfig struct {
	DataCenterId int64 `yaml:"data-center-id"`
	MachineId    int64 `yaml:"machine-id"`
	StartStmp    int64 `yaml:"start-stmp"`
}

type Configuration struct {
	Env struct {
		Active string `yaml:"active"`
	} `yaml:"env"`
	Gin struct {
		Release bool   `yaml:"release"`
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
	} `yaml:"gin"`
	Mysql struct {
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
		Database  string `yaml:"database"`
		DebugMode bool   `yaml:"debug-mode"`
	} `yaml:"mysql"`
	SnowFlake SnowFlakeConfig `yaml:"snow-flake"`
}

// default config instance
var Config IConfiguration

// Load config from application.yml, application-{env}.yml and cmd parameters.
// (`env` is given by application.yml)
func initConfig[T IConfiguration](config T) {
	// init by application.yml
	data, err := os.ReadFile("application.yml")
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}
	// init by application-{env}.yml
	if utf8.RuneCountInString(config.GetActiveEnv()) != 0 {
		data, err = os.ReadFile("application-" + config.GetActiveEnv() + ".yml")
		if err != nil {
			panic(err)
		}
		if err := yaml.Unmarshal(data, config); err != nil {
			panic(err)
		}
	}
	// set default values
	config.SolveDefaultValue()

	Config = config
	log.Println("config load complete")
}

func (config *Configuration) GetActiveEnv() string {
	return config.Env.Active
}

func (config *Configuration) GetGinRelease() bool {
	return config.Gin.Release
}

func (config *Configuration) GetGinAddr() string {
	return fmt.Sprintf("%s:%d", config.Gin.Host, config.Gin.Port)
}

func (config *Configuration) SolveDefaultValue() {
	if config.Gin.Port == 0 {
		config.Gin.Port = 5480
	}
	if config.SnowFlake.StartStmp == 0 {
		config.SnowFlake.StartStmp = 1626779686000
	}
}

func (config *Configuration) GetSnowFlakeConfig() SnowFlakeConfig {
	return config.SnowFlake
}
