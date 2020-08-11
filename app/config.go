package app

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config represents our App's configuration (secret-key, etc)
type Config struct {
	// A secret string used for session cookies, passwords, etc.
	SecretKey []byte
}

// InitConfig initializes our App's Config object based on viper or default values
// where none were provided if possible. Returns an error if a value was not provided
// and there is no default.
func InitConfig() (*Config, error) {
	config := &Config{
		SecretKey: []byte(viper.GetString("SecretKey")),
	}
	if len(config.SecretKey) == 0 {
		return nil, fmt.Errorf("SecretKey must be set")
	}
	return config, nil
}
