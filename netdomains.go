package main

import (
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

// NetDomain represents a broadcast domain.
// It holds metadata and a list of nodes.
type NetDomain struct {
	CfgDomain *ConfigDomain
	Nodes     []*NetNode
}

func (dom *NetDomain) fanoutBufferToNode(buf []byte, nn *NetNode) {
	nn.ConnLock.Lock()
	if nn.Conn == nil {
		nn.ConnLock.Unlock()
		return
	}
	if len(dom.CfgDomain.Jitter) == 2 {
		// Jitter is given as jitter[0] ms +/- jitter[1] ms
		sleepMs := dom.CfgDomain.Jitter[0] + (rand.Int31n(2*dom.CfgDomain.Jitter[1]) - dom.CfgDomain.Jitter[1])
		if Verbose {
			log.Println("Sleeping", sleepMs, "ms before delivering", len(buf), "bytes to", nn.Name)
		}
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)
	}
	_, err := nn.Conn.Write(buf)
	if err != nil {
		if Verbose {
			if err == io.EOF {
				log.Println("Node has disconnected:", nn.Name)
			} else {
				log.Println("Error reading connection:", nn.Name, err)
			}
		}
		nn.tearDownNode()
	}
	nn.ConnLock.Unlock()
}

// FanoutBuffer distributes the data in the given buf to all the nodes in this domain.
func (dom *NetDomain) fanoutBuffer(buf []byte, sender *NetNode) {
	for _, nn := range dom.Nodes {
		if nn == sender || nn.Conn == nil {
			continue
		}
		if rand.Float32() < dom.CfgDomain.Loss {
			if Verbose {
				log.Println("Lost", len(buf), "bytes in delivery to", nn.Name)
			}
			continue
		}
		go dom.fanoutBufferToNode(buf, nn)
	}
}

// NetNode represents one unique node in the network.
type NetNode struct {
	Name     string
	Domains  []*NetDomain
	Conn     net.Conn
	ConnLock sync.Mutex
}

// Run runs the node receiver loop.
func (nn *NetNode) Run() {
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
			nn.ConnLock.Lock()
			nn.tearDownNode()
			nn.ConnLock.Unlock()
			break
		}
		rbuf := make([]byte, rsize)
		copy(rbuf, buf)
		if Verbose {
			if rsize <= 16 {
				log.Println(nn.Name, "sent", rsize, "bytes:", rbuf, "fanning out to", len(nn.Domains), "domains")
			} else {
				log.Println(nn.Name, "sent, fanning out to", len(nn.Domains), "domains")
			}
		}
		for _, dom := range nn.Domains {
			dom.fanoutBuffer(rbuf, nn)
		}
	}
}

// tearDownNode is called when the node disconnects.
// It expects that the connection lock is held.
func (nn *NetNode) tearDownNode() {
	if nn.Conn == nil {
		return
	}
	if err := nn.Conn.Close(); err != nil {
		log.Printf("Connection close: %v", err)
	}
	nn.Conn = nil
}
