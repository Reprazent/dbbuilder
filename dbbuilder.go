package main

import (
	"fmt"
	// "strings"

	"flag"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"os"

	"path/filepath"
	// "regexp"
	// "time"
)

const Version = "0.0.1"

func red(in string) (out string) {
	in = "\033[31m" + in + "\033[0m"
	return in
}

func printError(message string) {
	fmt.Println(red(fmt.Sprintf("\n\t%s", message)))
}

// # Database Config struct.
//
// Each Rails `database.yml` environment block has the following structure, so
// it maps to this struct.
type DbConfig struct {
	Adapter  string
	Host     string
	Database string
	Username string
	Password string
}

// If no host is set, assume localhost.
func (self *DbConfig) SetDefaults() {
	if self.Host == "" {
		self.Host = "localhost"
	}
}

var currentDir = func() string {
	p, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return p
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s <environment>\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}

func getYamlPath(path string) (out string, err error) {
	if path == "" {
		path = currentDir()
	}
	if filepath.Ext(path) != ".yml" {
		path = filepath.Join(path, "config", "database.yml")
	}
	_, err = os.Stat(path)
	return path, err
}

func getConfig(yamlData []byte, env string) (*DbConfig, error) {
	m := make(map[string]*DbConfig)
	err := yaml.Unmarshal(yamlData, &m)
	if err != nil {
		return nil, err
	} // That did not go well
	config, ok := m[env]
	if !ok {
		var definedEnvs []string
		for key, _ := range m {
			definedEnvs = append(definedEnvs, key)
			printError(fmt.Sprintf("Environment %s not found, use one of %s", env, definedEnvs))
		}
		return nil, err
	}
	config.SetDefaults()
	return config, nil
}

func createUserCommand(config *DbConfig) string {
	return fmt.Sprintf("createuser -h %s -d -R -e -w %s", config.Host, config.Username)
}

func createDatabaseCommand(config *DbConfig) string {
	return fmt.Sprintf("createdb -h %s -O  %s -U %s -d -w -e %s", config.Host, config.Username, config.Username, config.Database)
}

func main() {
	flag.Usage = usage
	var version bool
	var path string
	var environment string
	flag.BoolVar(&version, "v", false, "Prints current version")
	flag.StringVar(&path, "p", "config/database.yml", "Path to yaml (otherwise config/database.yml)")
	flag.StringVar(&environment, "e", "test", "Set the database test env to create")
	flag.Parse()

	if version {
		fmt.Println(Version)
		os.Exit(0)
	}

	filePath, err := getYamlPath(path)

	if err != nil {
		printError("File does not exist")
		os.Exit(2)
	}

	fmt.Println("Using file configuration file:", filePath)

	yamlData, err := ioutil.ReadFile(filePath)

	if err != nil {
		printError(fmt.Sprintf("Could not read yaml file %v", err))
		os.Exit(2)
	}

	config, err := getConfig(yamlData, environment)
	if err != nil {
		os.Exit(2)
	}

	fmt.Println(createUserCommand(config))
	fmt.Println(createDatabaseCommand(config))


}
