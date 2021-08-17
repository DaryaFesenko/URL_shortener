package app

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Postgres string `yaml:"postgres"`
}

func NewConfig(fileName string) (*Config, error) {
	var data []byte
	var conf Config

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return &conf, err
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return &conf, err
	}

	return &conf, nil
}

func (c *Config) GetPostgres() string {
	return c.Postgres
}
