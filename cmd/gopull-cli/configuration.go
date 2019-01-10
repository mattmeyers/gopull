package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// InitConfig initializes the viper configuration.
func InitConfig(filename string) {
	setDefaults()

	viper.SetConfigName(filename)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %s", err))
	}
}

func setDefaults() {
	viper.SetDefault("repos_dir", "$HOME/repos")
	viper.SetDefault("gopull_dir", "$GOPATH/src/github.com/mattmeyers/gopull")
}
