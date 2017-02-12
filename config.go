package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	StaticDir   string `json:"static_dir"`
	TemplateDir string `json:"template_dir"`
	PageDir     string `json:"page_dir"`
	OutputDir   string `json:"output_dir"`
}

func (c *Config) Verify() error {
	return nil
}

func loadConfig(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open config file '%s': %v", file, err)
	}

	cfg := &Config{}
	dec := json.NewDecoder(f)
	if err := dec.Decode(cfg); err != nil {
		return cfg, fmt.Errorf("Couldn't parse config file '%s': %v", file, err)
	}
	return cfg, nil
}
