package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is a structure containing configuration fields for this application.
type Config struct {
	Loglevel         string `env:"LOG_LEVEL"     env-default:"error"`
	HTTPListenIPPort string `env:"HTTP_LISTEN"   env-default:":80"`
	JWT_SECRET       string `env:"JWT_SECRET"`    // stored in secret
	PASSWORD_SALT    string `env:"PASSWORD_SALT"` // stored in secret
	DB_HOST          string `env:"DB_HOST"       env-default:"127.0.0.1"`
	DB_PORT          string `env:"DB_PORT"       env-default:"5432"`
	DB_USER          string `env:"DB_USER"       env-default:"go-api"`
	DB_PASS          string `env:"DB_PASS"` // stored in secret
	DB_DATABASE      string `env:"DB_DATABASE"   env-default:"go-api"`
}

var Cfg *Config

func init() {
	Cfg = &Config{}
	err := cleanenv.ReadEnv(Cfg)
	if err != nil {
		fmt.Printf("Something went wrong while reading the configuration: %s", err)
		os.Exit(1)
	}
}
