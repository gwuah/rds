package config

import (
	"bytes"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

func ParseConfig(r io.Reader, config *Config) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	data = []byte(os.ExpandEnv(string(data)))
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	return dec.Decode(&config)
}

func ParseConfigFromFile(filename string) (*Config, error) {
	var cfg *Config
	f, err := os.Open(filename)
	if err != nil {
		return cfg, err
	}
	defer func() { _ = f.Close() }()
	return cfg, ParseConfig(f, cfg)
}
