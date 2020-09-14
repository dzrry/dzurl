package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}

type Config struct {
	Redis  *RedisConfig  `yaml:"redis"`
	Server *ServerConfig `yaml:"server"`
}

func Read(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(file, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
