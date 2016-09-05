package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

// Cfg is the global configuration data.
var Cfg ConfigMain

// AllNodes is the global list of all nodes.
var AllNodes map[string]*NetNode

// AllDomains is the global list of all domains.
var AllDomains []NetDomain

// Verbose is the global log level.
var Verbose bool

func showUsage() {
	fmt.Println("usage:", os.Args[0], "[-c config.json] [-v]")
}

// handleConnection handles the config phase of the connection.
func handleConnection(conn net.Conn) {
	var buf = make([]byte, 1)
	var bline []byte
	for {
		r, err := conn.Read(buf)
		if err != nil {
			log.Fatalln("Error reading config line from:", conn)
			return
		}
		if r != 1 {
			return
		}
		if buf[0] == '\n' {
			line := string(bline)
			if nn, ok := AllNodes[line]; ok {
				if nn.Conn != nil {
					log.Println("Node", line, "already connected! Refusing another connection.")
					return
				}
				nn.Conn = conn
				if Verbose {
					log.Println(conn, "is", line)
				}
				nn.Run()
			} else {
				log.Fatalln("First line in new connection needs to be the node name for:", conn)
			}
			return
		}
		bline = append(bline, buf[0])
	}
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "config.json", "Config file name")
	flag.StringVar(&configFile, "c", "config.json", "Short for config file name")
	flag.BoolVar(&Verbose, "v", false, "Verbose output")
	flag.Parse()

	if _, err := os.Stat(configFile); err != nil {
		fmt.Printf("Config file not found: '%s'\n", configFile)
		showUsage()
		flag.PrintDefaults()
		return
	}

	var err error
	Cfg, err = ReadConfig(configFile)
	if err != nil {
		fmt.Printf("Cannot parse config file: '%s': %v\n", configFile, err)
		return
	}

	if Verbose {
		fmt.Println(Cfg)
	}

	checkConfig()

	// Parse the config into NetNode and NetDomain slices
	AllNodes = generateAllNodes(&Cfg)
	AllDomains = generateAllDomains(&Cfg)

	if Verbose {
		log.Println("Working with", len(AllDomains), "domains and", len(AllNodes), "nodes.")
	}

	srv, err := net.Listen("tcp", fmt.Sprintf(":%d", Cfg.Port))
	if err != nil {
		log.Fatalln("Error creating a TCP listener:", err)
	}

	// Main accept() loop
	for {
		conn, err := srv.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			break
		}
		go handleConnection(conn)
	}
}

func generateAllNodes(cfg *ConfigMain) map[string]*NetNode {
	all := make(map[string]*NetNode)
	for di := range cfg.Network {
		for _, nname := range cfg.Network[di].Nodes {
			if _, ok := all[nname]; !ok {
				all[nname] = &NetNode{Name: nname}
			}
		}
	}
	return all
}

func generateAllDomains(cfg *ConfigMain) []NetDomain {
	all := make([]NetDomain, len(cfg.Network))
	for di := range cfg.Network {
		all[di].CfgDomain = &cfg.Network[di]
		all[di].Nodes = make([]*NetNode, len(cfg.Network[di].Nodes))
		for ni, nname := range cfg.Network[di].Nodes {
			all[di].Nodes[ni] = AllNodes[nname]
			AllNodes[nname].Domains = append(AllNodes[nname].Domains, &all[di])
		}
	}
	return all
}
