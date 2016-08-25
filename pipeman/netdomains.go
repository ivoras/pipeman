package main

import (
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
)

// NetDomain represents a broadcast domain. It holds metadata and a list of nodes.
type NetDomain struct {
	CfgDomain *PipemanDomain
	Nodes     []*NetNode
}

// FanoutBuffer distributes the data in the given buf to the nodes in this domain
func (dom *NetDomain) FanoutBuffer(buf []byte, sender *NetNode) {
	for _, nn := range dom.Nodes {
		if nn == sender {
			continue
		}
		if rand.Float32() < dom.CfgDomain.Loss {
			continue
		}
		nn.ConnLock.Lock()
		if nn.Conn == nil {
			nn.ConnLock.Unlock()
			continue
		}
		_, err := nn.Conn.Write(buf)
		if err != nil {
			if Verbose {
				if err == io.EOF {
					log.Println("Node has disconnected:", nn.Name)
				} else {
					log.Println("Error reading connection", nn.Name, err)
				}
			}
			TearDownNode(nn)
		}
		nn.ConnLock.Unlock()
	}
}

// NetNode represents one unique node in the network
type NetNode struct {
	Name     string
	Domains  []*NetDomain
	Conn     *net.TCPConn
	ConnLock sync.Mutex
}

// NetNodeRun runs the node receiver loop
func (nn *NetNode) NetNodeRun() {
	buf := make([]byte, Cfg.BufferSize)
	for {
		nn.ConnLock.Lock()
		if nn.Conn == nil {
			nn.ConnLock.Unlock()
			break
		}
		rsize, err := nn.Conn.Read(buf)
		nn.ConnLock.Unlock()
		if err != nil {
			if Verbose {
				if err == io.EOF {
					log.Println("Node has disconnected:", nn.Name)
				} else {
					log.Println("Error reading connection", nn.Name, err)
				}
			}
			TearDownNode(nn)
			break
		}
		if rsize == 0 {
			continue
		}
		if Verbose {
			log.Println(nn.Name, "sent", rsize, "bytes", "; fanning to", len(nn.Domains), "domains")
		}
		rbuf := buf[:rsize]
		for _, dom := range nn.Domains {
			dom.FanoutBuffer(rbuf, nn)
		}
	}
}
