package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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

func createUserCommand(config *DbConfig) *exec.Cmd {
	return exec.Command("createuser", "-h", config.Host, "-S", "-d", "-R", "-e", "-w", config.Username)
}

func createDatabaseCommand(config *DbConfig) *exec.Cmd {
	return exec.Command("createdb", "-h", config.Host, "-O", config.Username, "-U", config.Username, "-d", "-w", "-e", config.Database)
}

func runCommand(cmd *exec.Cmd) error {
	fmt.Println("Running: ", cmd)
	output, err := cmd.Output()
	fmt.Println(string(output))
	return err
}

func handleError(err error, clean bool) {
	if err != nil {
		printError(err.Error())
		if clean != true {
			printError("Exiting...")
			os.Exit(1)
		}
	}
}

func main() {
	flag.Usage = usage
	var version bool
	var path string
	var environment string
	var clean bool
	flag.BoolVar(&version, "v", false, "Prints current version")
	flag.BoolVar(&clean, "c", false, "Don't exit on errors, useful for CI")
	flag.StringVar(&path, "p", "config/database.yml", "Path to yaml (otherwise config/database.yml)")
	flag.StringVar(&environment, "e", "test", "Set the database test env to create")
	flag.Parse()

	if version {
		fmt.Println(Version)
		os.Exit(0)
	}

	filePath, err := getYamlPath(path)

	handleError(err, clean)

	fmt.Println("Using file configuration file:", filePath)

	yamlData, err := ioutil.ReadFile(filePath)

	handleError(err, clean)

	config, err := getConfig(yamlData, environment)

	handleError(err, clean)

	err = runCommand(createUserCommand(config))

	handleError(err, clean)

	err = runCommand(createDatabaseCommand(config))

	handleError(err, clean)

	fmt.Println("Done...")

	os.Exit(0)
}
