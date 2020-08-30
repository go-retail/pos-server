package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func getConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("conf")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			failOnError(err, "Error reading Config file")
		}
	}
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())
}
