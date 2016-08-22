package main

import (
	"flag"
	"fmt"
	"os"
)

var cfg PipemanConfig

func showUsage() {
	fmt.Println("usage:", os.Args[0], "[-c config.json]")
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "config.json", "Config file name")
	flag.StringVar(&configFile, "c", "config.json", "Short for config file name")
	flag.Parse()

	if _, err := os.Stat(configFile); err != nil {
		fmt.Printf("Config file not found: '%s'\n", configFile)
		showUsage()
		flag.PrintDefaults()
		return
	}

	fmt.Println(configFile)
	cfg, err := ReadConfig(configFile)
	if err != nil {
		fmt.Printf("Cannot parse config file: '%s': %v\n", configFile, err)
		return
	}

	fmt.Println(cfg)
}
