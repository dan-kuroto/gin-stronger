package gs

import (
	"log"
	"os"
	"unicode/utf8"

	"github.com/dan-kuroto/gin-stronger/config"
	"gopkg.in/yaml.v3"
)

// default config instance
var Config config.IConfiguration

// Load config from application.yml, application-{env}.yml and cmd parameters.
// (`env` is given by application.yml)
func initConfig[T config.IConfiguration](config T) {
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
