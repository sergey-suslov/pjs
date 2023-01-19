package config

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

const CONFIG_DIR_NAME = ".pje"
const CONFIG_ROOT_FILE_NAME = ".config.toml"

var BadConfig = errors.New("Error reading config")
var NoConfigFound = errors.New("No config found")

func TryParseLocalConfig() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(path.Join(cwd, CONFIG_DIR_NAME))
	if os.IsNotExist(err) {
		return nil, NoConfigFound
	}
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, BadConfig
	}

	configData, err := os.ReadFile(path.Join(cwd, CONFIG_DIR_NAME, CONFIG_ROOT_FILE_NAME))
	if err != nil {
		return nil, NoConfigFound
	}

	var conf Config
	_, err = toml.Decode(string(configData[:]), &conf)
	if err != nil {
		return nil, BadConfig
	}
	return &conf, nil
}

func TryCreateLocalConfig(projectId, projectName string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	configDirPath := path.Join(cwd, CONFIG_DIR_NAME)
	_, err = os.Stat(configDirPath)
	if os.IsNotExist(err) {
		if err := os.Mkdir(configDirPath, 0744); err != nil {
			return err
		}
	}
	var configBuffer bytes.Buffer
	err = toml.NewEncoder(&configBuffer).Encode(Config{ProjectName: projectName, ProjectId: projectId})
	if err != nil {
		return err
	}
	configBytes, err := ioutil.ReadAll(&configBuffer)
	err = os.WriteFile(path.Join(configDirPath, CONFIG_ROOT_FILE_NAME), configBytes, 0744)
	if err != nil {
		return err
	}
	return nil
}
