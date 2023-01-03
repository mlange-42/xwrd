package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mlange-42/xwrd/util"
	"gopkg.in/yaml.v3"
)

var (
	// ErrNoConfig is an error for no config file available
	ErrNoConfig = errors.New("no config file")
)

// Config for track
type Config struct {
	Dict string `yaml:"dict"`
}

// GetDict return the current dictionary
func (c *Config) GetDict() util.Dict {
	return util.NewDict(c.Dict)
}

// LoadConfig loads the track config, or creates and saves default settings
func LoadConfig() (Config, error) {
	conf, err := tryLoadConfig()
	if err == nil {
		return conf, nil
	}
	if !errors.Is(err, ErrNoConfig) {
		return conf, err
	}

	conf = Config{
		Dict: "en/yawl",
	}

	err = SaveConfig(conf)
	if err != nil {
		return Config{}, fmt.Errorf("could not save config file: %s", err)
	}

	return conf, nil
}

func tryLoadConfig() (Config, error) {
	file, err := ioutil.ReadFile(util.ConfigPath())
	if err != nil {
		return Config{}, ErrNoConfig
	}

	var conf Config

	if err := yaml.Unmarshal(file, &conf); err != nil {
		return Config{}, err
	}

	if err = CheckConfig(&conf); err != nil {
		return conf, err
	}

	return conf, nil
}

// SaveConfig saves the given config to it's default location
func SaveConfig(conf Config) error {
	if err := CheckConfig(&conf); err != nil {
		return err
	}

	path := util.ConfigPath()

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	bytes, err := yaml.Marshal(&conf)
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}

// CheckConfig checks a config for consistency
func CheckConfig(conf *Config) error {
	return nil
}
