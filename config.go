package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
	err = json.Unmarshal(f, &cfg)
	return cfg, err
}

func checkConfig() {
	if Cfg.Type != "pipeman" {
		log.Fatalln("Invalid config file (missing type:\"pipeman\")")
	}
	if Cfg.BufferSize < 1 {
		log.Fatalln("buffer_size must be at least 1")
	}
	if Cfg.Port < 0 {
		log.Fatalln("port must be a positive integer")
	}
	checkConfigDomain()
}

func checkConfigDomain() {
	for _, pd := range Cfg.Network {
		if pd.Loss < 0 {
			log.Fatalln("Lost must be a positive decimal number")
		}
		if len(pd.Jitter) != 0 {
			if len(pd.Jitter) != 2 {
				log.Fatalln("Jitter must be specified as an array of 2 numbers: (#num1 +/- #num2) milliseconds")
			}
			for _, j := range pd.Jitter {
				if j < 0 {
					log.Fatalln("Jitter spec must be positive integers")
				}
			}
			if pd.Jitter[1] > pd.Jitter[0] {
				log.Fatalln("Jitter must be specified as an array of 2 numbers: (#num1 +/- #num2) milliseconds, #num2 < #num1")
			}
		}
	}
}
