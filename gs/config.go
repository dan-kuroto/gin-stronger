package gs

import (
	"log"

	"github.com/dan-kuroto/gin-stronger/config"
)

// default config instance
var Config config.IConfiguration

// Load config from application.yml, application-{env}.yml and cmd parameters.
// (`env` is given by application.yml)
func InitConfig[T config.IConfiguration](config T) error {
	Config = config

	const baseName = "application"
	// init by application.yml
	if err := config.InitByYaml(baseName, ""); err != nil {
		return err
	}
	// init by application-{env}.yml
	if config.GetActiveEnv() != "" {
		if err := config.InitByYaml(baseName, config.GetActiveEnv()); err != nil {
			return err
		}
	}
	// set default values
	config.SolveDefaultValue()

	log.Println("config load complete")
	return nil
}
