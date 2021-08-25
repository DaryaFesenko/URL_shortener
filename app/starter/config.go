package starter

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Postgres       string        `yaml:"postgres"`
	HashSalt       string        `yaml:"hashSalt"`
	SigningKey     []byte        `yaml:"signingKey"`
	ExpireDuration time.Duration `yaml:"expireDuration"`
}

func NewConfig(fileName string) (*Config, error) {
	var data []byte
	var conf Config

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("can't read config file: %v", err)
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal config file: %v", err)
	}

	return &conf, nil
}

func (c *Config) GetPostgres() string {
	return c.Postgres
}
