package api

import "github.com/spf13/viper"

// Config is an object representing the configuration of our API (port, proxy-count, etc.)
type Config struct {
	// The port to bind the web application server to
	Port int

	// The number of proxies positioned in front of the API.
	// This is used to interpret X-Forwarded-For headers.
	ProxyCount int
}

// InitConfig initializes our API's Config object using viper and setting defaults
// where values are not provided.
func InitConfig() (*Config, error) {
	config := &Config{
		Port:       viper.GetInt("Port"),
		ProxyCount: viper.GetInt("ProxyCount"),
	}
	if config.Port == 0 {
		config.Port = 9092
	}
	return config, nil
}
