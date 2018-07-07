package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseHost string `envconfig:"ls_db_host" default:"localhost"`
	DatabasePort string `envconfig:"ls_db_port" default:"3306"`
	DatabaseName string `envconfig:"ls_db_name" default:"livestock"`
	DatabaseUser string `envconfig:"ls_db_user" default:"root"`
	DatabasePass string `envconfig:"ls_db_pass" default:""`
	DebugMode    bool   `envconfig:"ls_debug" default:"true"`
}

func (cfg *Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"%s:%s@%s:%s/%s?parseTime=true",
		cfg.DatabaseUser,
		cfg.DatabasePass,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	)
}

func Get() Config {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		panic("Failed to read configuration from environment variables")
	}
	return cfg
}
