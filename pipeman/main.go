package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

// Cfg is the global configuration data
var Cfg PipemanConfig

// AllNodes is the global list of all nodes
var AllNodes map[string]*NetNode

// AllDomains is the global list of all domains
var AllDomains []NetDomain

// Verbose is the global log level
var Verbose bool

func showUsage() {
	fmt.Println("usage:", os.Args[0], "[-c config.json] [-v]")
}

// Handles the config phase of the connection
func handleConnection(conn *net.TCPConn) {
	var buf = make([]byte, 1)
	var bline []byte
	for {
		r, err := conn.Read(buf)
		if err != nil {
			log.Fatalln("Error reading config line from:", *conn)
			return
		}
		if r != 1 {
			return
		}
		if buf[0] == '\n' {
			line := string(bline)
			if nn, ok := AllNodes[line]; ok {
				nn.Conn = conn
				if Verbose {
					log.Println(conn, "is", line)
				}
				nn.NetNodeRun()
			} else {
				log.Fatalln("First line in config state needs to be the node name for:", *conn)
			}
			return
		}
		bline = append(bline, buf[0])
	}
}

// tearDownNode is called when the node disconnects. It expects that the connection lock is held.
func tearDownNode(nn *NetNode) {
	if nn.Conn != nil {
		nn.Conn.Close()
		nn.Conn = nil
	}
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

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "config.json", "Config file name")
	flag.StringVar(&configFile, "c", "config.json", "Short for config file name")
	flag.BoolVar(&Verbose, "v", false, "Verbose output")
	flag.Parse()

	var err error

	if _, err := os.Stat(configFile); err != nil {
		fmt.Printf("Config file not found: '%s'\n", configFile)
		showUsage()
		flag.PrintDefaults()
		return
	}

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
	AllNodes = make(map[string]*NetNode)
	for di := range Cfg.Network {
		for _, nname := range Cfg.Network[di].Nodes {
			if _, ok := AllNodes[nname]; ok {
				// Already exists
				continue
			}
			AllNodes[nname] = new(NetNode)
			AllNodes[nname].Name = nname
		}
	}

	AllDomains = make([]NetDomain, len(Cfg.Network))
	for di := range Cfg.Network {
		AllDomains[di].CfgDomain = &Cfg.Network[di]
		AllDomains[di].Nodes = make([]*NetNode, len(Cfg.Network[di].Nodes))
		for ni, nname := range Cfg.Network[di].Nodes {
			AllDomains[di].Nodes[ni] = AllNodes[nname]
			AllNodes[nname].Domains = append(AllNodes[nname].Domains, &AllDomains[di])
		}
	}

	addr := net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: Cfg.Port}
	srv, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		log.Fatalln("Error creating a TCP listener:", err)
	}

	for {
		conn, err := srv.AcceptTCP()
		if err != nil {
			log.Println("Error accepting connection:", err)
			break
		}
		if Verbose {
			log.Println(conn)
		}
		go handleConnection(conn)
	}
}
