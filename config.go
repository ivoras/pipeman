package main

import (
	"encoding/json"
	"io/ioutil"
)

// ConfigMain holds the JSON configuration structure for Pipeman
type ConfigMain struct {
	Type       string         `json:"type"`
	Port       int            `json:"port"`
	BufferSize uint32         `json:"buffer_size"`
	Network    []ConfigDomain `json:"network"`
}

// ConfigDomain holds a list of nodes belonging to the same group, with some metadata
type ConfigDomain struct {
	Name   string   `json:"name"`
	Nodes  []string `json:"nodes"`
	Loss   float32  `json:"loss"`
	Jitter []int32  `json:"jitter"`
}

// ReadConfig reads a config file and returns the parsed PipemanConfig struct
func ReadConfig(fileName string) (ConfigMain, error) {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		return ConfigMain{}, err
	}

	var cfg ConfigMain
	if err = json.Unmarshal(f, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
