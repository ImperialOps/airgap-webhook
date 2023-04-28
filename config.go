package main

import (
	"errors"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	cfgFile    string    `json:"cfgFile"`
	listenAddr string    `json:"listenAddr"`
	apiAddr    string    `json:"apiAddr"`
	tls        ConfigTls `json:"tls"`
}

type ConfigTls struct {
	enabled  bool   `json:"enabled"`
	certFile string `json:"certFile"`
	keyFile  string `json:"keyFile"`
}

func NewConfig() (*Config, error) {
	config := Config{
		cfgFile:    "",
		listenAddr: "0.0.0.0:8080",
		apiAddr:    "",
		tls: ConfigTls{
			enabled:  false,
			certFile: "",
			keyFile:  "",
		},
	}

	pflag.StringVar(&config.cfgFile, "config", config.cfgFile, "config file location")
	pflag.StringVar(&config.listenAddr, "listen-address", config.listenAddr, "server listen address")
	pflag.BoolVar(&config.tls.enabled, "tls-enabled", config.tls.enabled, "controls whether tls is enabled, good for testing")
	pflag.StringVar(&config.tls.certFile, "tls-cert", config.tls.certFile, "tls certificate to serve")
	pflag.StringVar(&config.tls.keyFile, "tls-key", config.tls.keyFile, "tls key")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if config.cfgFile != "" {
		viper.SetConfigFile(config.cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			return &config, err
		}
	}

	viper.SetEnvPrefix("ag")
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&config); err != nil {
		return &config, err
	}

	// validate
	if config.tls.enabled {
		if config.tls.certFile == "" {
			return &config, errors.New("must supply certificate file")
		}
		if config.tls.keyFile == "" {
			return &config, errors.New("must supply private key file")
		}
	}

	return &config, nil
}
