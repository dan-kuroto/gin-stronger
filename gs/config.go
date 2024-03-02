package gs

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/dan-kuroto/gin-stronger/config"
)

// default config instance
var Config config.IConfiguration

// Load config from application.yml, application-{env}.yml and cmd parameters.
// (`env` is given by application.yml)
func InitConfig[T config.IConfiguration](config T) error {
	Config = config
	// set default values
	defer config.SolveDefaultValue()

	const baseName = "application"
	// init by application.yml
	if err := initConfigByYaml(config, baseName, ""); err != nil {
		return err
	}
	// init by application-{env}.yml
	if config.GetActiveEnv() != "" {
		if err := initConfigByYaml(config, baseName, config.GetActiveEnv()); err != nil {
			return err
		}
	}
	log.Println("config load complete")
	return nil
}

// If env is empty string, init config by {name}.yml, otherwise {name}-{env}.yml
func initConfigByYaml[T config.IConfiguration](config T, baseName string, env string) error {
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
