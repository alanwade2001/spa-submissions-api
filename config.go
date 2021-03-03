package main

import "github.com/spf13/viper"

// ConfigService s
type ConfigService struct {
}

// Load f
func (cs ConfigService) Load() error {

	//viper.AddConfigPath(".")
	//viper.SetConfigName("app")
	//viper.SetConfigType("env")

	viper.AutomaticEnv()

	//err := viper.ReadInConfig()
	//if err != nil {
	//	panic(err)
	//}

	return nil
}

// NewConfigService f
func NewConfigService() ConfigAPI {
	return &ConfigService{}
}
