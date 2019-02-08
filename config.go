package gopull

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/spf13/viper"
)

// configPaths is a slice of possible paths to the config files.
var configPaths = []string{"$GOPATH/src/github.com/mattmeyers/gopull", getCallerDir(), "."}

// NewConfig initializes a local repo configuration viper instance.
//
// Viper looks in three places for the config.json configuration file.
// First it looks in the standard path that go get would place the source
// code. Then it would look in directory containing the source code for
// the binary calling this function. Finally it will look in the current
// working directory.
func NewConfig() {
	setDefaults()

	viper.SetConfigName("config")
	viper.AddConfigPath(configPaths[0])
	viper.AddConfigPath(configPaths[1])
	viper.AddConfigPath(configPaths[2])
	err := viper.ReadInConfig()
	if err != nil {
		var file *os.File
		for _, path := range configPaths {
			file, err = os.Create(fmt.Sprintf("%s/config.json", path))
			if err == nil {
				break
			}
		}

		if file == nil {
			panic(fmt.Errorf("fatal error creating config.json: %s", err))
		}
		defer file.Close()

		if err = viper.WriteConfig(); err != nil {
			panic(fmt.Errorf("fatal error writing default config.json: %s", err))
		}
	}
}

func getCallerDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("No caller information.")
	}
	return path.Dir(filename)
}

func setDefaults() {
	viper.SetDefault("paths", map[string]string{
		"repos_dir":   "$HOME/repos",
		"gopull_dir":  "$GOPATH/src/github.com/mattmeyers/gopull/cmd/gopull-api",
		"scripts_dir": "$GOPATH/src/github.com/mattmeyers/gopull/cmd/gopull-api/deployment_scripts",
	})
	viper.SetDefault("repos", map[string]map[string]string{})
}

func createNewConfig(filename string) (*os.File, error) {
	var file *os.File
	var err error

	for _, path := range configPaths {
		file, err = os.Create(fmt.Sprintf("%s/%s", path, filename))
		if err == nil {
			break
		}
	}

	if file == nil {
		return file, fmt.Errorf("error creating %s", filename)
	}
	return file, nil
}
