package env

import (
	"fmt"

	"github.com/spf13/viper"
)

// loads the enviroment variable and set this
func InitEnv(serviceName string, mandatoryParams []string, optionalParams map[string]interface{}) error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(serviceName)

	notPresent := make([]string, 0)
	for _, i := range mandatoryParams {
		if !viper.IsSet(i) {
			notPresent = append(notPresent, i)
		}
	}

	if len(notPresent) != 0 {
		return fmt.Errorf("mandatory env missing %+v", notPresent)
	}

	for key, val := range optionalParams {
		viper.SetDefault(key, val)
	}

	return nil
}
