package config

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Address string `yaml:"address"`
}

type Config struct {
	ServerConfig ServerConfig `yaml:"server_config"`
}

func ParseConfig(fp string) (*Config, error) {
	var config Config
	f, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(f, &config); err != nil {
		panic(err)
	}
	return &config, nil
}

func ProvideConfig(l *zap.Logger) *Config {
	conf, err := ParseConfig("./config.yaml")
	if err != nil {
		l.Fatal("cannot parse config", zap.Error(err))
	}
	return conf
}
