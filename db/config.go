package db

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config represents our App's database configuration (for now only the DatabaseURI)
type Config struct {
	DatabaseURI string
}

// InitConfig initializes our App's Database configuration object using values obtained from viper,
// or defaults where possible. If a value is not provided and there is no default value it returns
// an error.
func InitConfig() (*Config, error) {
	config := &Config{
		DatabaseURI: viper.GetString("DatabaseURI"),
	}
	if config.DatabaseURI == "" {
		return nil, fmt.Errorf("DatabaseURI must be set")
	}
	return config, nil
}
