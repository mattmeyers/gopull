package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/spf13/viper"
)

// InitConfig initializes the viper configuration.
//
// The configuration is loaded from the config.json file found in
// the gopull-cli directory. If this file is not present or cannot
// be accessed, it is created and the default values are written to
// it.
func InitConfig() {
	setDefaults()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("No caller information.")
	}
	dir := path.Dir(filename)

	viper.SetConfigName("config")
	viper.AddConfigPath(dir)
	err := viper.ReadInConfig()
	if err != nil {
		if _, err = os.Create(fmt.Sprintf("%s/config.json", dir)); err != nil {
			panic(fmt.Errorf("fatal error creating config.json: %s", err))
		}
		if err = viper.WriteConfig(); err != nil {
			panic(fmt.Errorf("fatal error writing default config: %s", err))
		}

	}
}

func setDefaults() {
	viper.SetDefault("repos_dir", "$HOME/repos")
	viper.SetDefault("gopull_dir", "$GOPATH/src/github.com/mattmeyers/gopull")
}
