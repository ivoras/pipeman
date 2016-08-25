package main

import (
	"log"
	"math/rand"
	"net"
)

// NetDomain represents a broadcast domain. It holds metadata and a list of nodes.
type NetDomain struct {
	CfgDomain *PipemanDomain
	Nodes     []*NetNode
}

// FanoutBuffer distributes the data in the given buf to the nodes in this domain
func (dom *NetDomain) FanoutBuffer(buf []byte, sender *NetNode) {
	for _, nn := range dom.Nodes {
		if nn == sender || nn.Conn == nil {
			continue
		}
		if rand.Float32() < dom.CfgDomain.Loss {
			continue
		}
		nn.Conn.Write(buf)
	}
}

// NetNode represents one unique node in the network
type NetNode struct {
	Name    string
	Domains []*NetDomain
	Conn    *net.TCPConn
}

// NetNodeRun runs the node receiver loop
func (nn *NetNode) NetNodeRun() {
	var buf = make([]byte, Cfg.BufferSize)
	for {
		rsize, err := nn.Conn.Read(buf)
		if err != nil {
			log.Fatalln("Error reading connection", nn.Conn)
			break
		}
		rbuf := buf[:rsize]
		for _, dom := range nn.Domains {
			dom.FanoutBuffer(rbuf, nn)
		}
	}
}
