package gs

import (
	"flag"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

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
	SnowFlake struct {
		DataCenterId int64 `yaml:"data-center-id"`
		MachineId    int64 `yaml:"machine-id"`
	} `yaml:"snow-flake"`
}

// default config instance
var Config Configuration

// It is shorthand for gs.Init(&gs.Config)
func InitDefault() {
	Init(&Config)
	log.Println("config load complete")
}

// Load config from application.yml, application-{env}.yml and cmd parameters.
// (`env` is given by applicaiotn.yml)
func Init(config *Configuration) {
	// init by application.yml
	data, err := os.ReadFile("application.yml")
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}
	// init by application-{env}.yml
	if utf8.RuneCountInString(config.Env.Active) != 0 {
		data, err = os.ReadFile("application-" + config.Env.Active + ".yml")
		if err != nil {
			panic(err)
		}
		if err := yaml.Unmarshal(data, config); err != nil {
			panic(err)
		}
	}
	// init by cmd parameters
	release := flag.Bool("release", false, "use release mode")
	host := flag.String("host", "", "gin host")
	port := flag.Int("port", 0, "gin port")
	flag.Parse()
	if *release {
		config.Gin.Release = true
	}
	if *host != "" {
		config.Gin.Host = *host
	}
	if *port != 0 {
		config.Gin.Port = *port
	}
}

func (config *Configuration) GetGinAddr() string {
	if config.Gin.Port == 0 {
		return fmt.Sprintf("%s:%d", config.Gin.Host, 8080)
	} else {
		return fmt.Sprintf("%s:%d", config.Gin.Host, config.Gin.Port)
	}
}
