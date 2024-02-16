package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type IConfiguration interface {
	GetActiveEnv() string
	GetGinRelease() bool
	GetGinAddr() string
	GetSnowFlakeConfig() SnowFlakeConfig

	InitByYaml(baseName string, env string) error
	SolveDefaultValue()
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
	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
	SnowFlake SnowFlakeConfig `yaml:"snow-flake"`
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

// If env is empty string, init config by {name}.yml, otherwise {name}-{env}.yml
func (config *Configuration) InitByYaml(baseName string, env string) error {
	var fpath string
	if env == "" {
		fpath = baseName + ".yml"
	} else {
		fpath = baseName + "-" + env + ".yml"
	}

	data, err := os.ReadFile(fpath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}
	return nil
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
