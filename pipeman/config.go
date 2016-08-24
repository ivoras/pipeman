package main

import (
	"encoding/json"
	"io/ioutil"
)

// PipemanConfig holds the JSON configuration structure for Pipeman
type PipemanConfig struct {
	Type       string         `json:"type"`
	Port       int            `json:"port"`
	BufferSize int            `json:"buffer_size"`
	Network    []PipemanDomain `json:"network"`
}

// PipemanGroup holds a list of nodes belonging to the same group, with some metadata
type PipemanDomain struct {
	Name  string   `json:"name"`
	Nodes []string `json:"nodes"`
	Loss  float32  `json:"loss"`
}

// ReadConfig reads a config file and returns the parsed PipemanConfig struct
func ReadConfig(fileName string) (PipemanConfig, error) {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		return PipemanConfig{}, err
	}

	var cfg PipemanConfig
	if err = json.Unmarshal(f, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
