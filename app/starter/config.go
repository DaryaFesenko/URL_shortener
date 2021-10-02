package starter

import (
	"os"
	"time"
)

type Config struct {
	Postgres       string        `yaml:"postgres"`
	HashSalt       string        `yaml:"hashSalt"`
	SigningKey     string        `yaml:"signingKey"`
	ExpireDuration time.Duration `yaml:"expireDuration"`
	ServerAddress  string        `yaml:"serverAddress"`
}

func NewConfig() (*Config, error) {
	var conf Config

	conf.Postgres = os.Getenv("PG_DSN")
	conf.ServerAddress = os.Getenv("PORT")
	conf.HashSalt = "hashSalt"
	conf.SigningKey = "signingKey"
	conf.ExpireDuration = time.Hour * 24 * 7

	return &conf, nil
}
