package main

import (
	"log"
	"net"
	"time"
)

// refer :https://github.com/jyootai/DhtBT/

var BOOTSTRAP []string = []string{
	"67.215.246.10:6881",  //router.bittorrent.com
	"91.121.59.153:6881",  //dht.transmissionbt.com
	"82.221.103.244:6881", //router.utorrent.com
	"212.129.33.50:6881"}

type Id []byte

type NodeInfo struct {
	ip       net.IP
	port     int
	id       Id
	lastSeen time.Time
}

type Routing struct {
	selfNode *KNode
	table    []*Bucket
}

type Bucket struct {
	nodes      []*NodeInfo
	lastChange time.Time
}

type Network struct {
	dhtNode *KNode
	conn    *net.Conn
}

type KNode struct {
	node    *NodeInfo
	routing *Routing
	network *Network
	log     *log.Logger
	krpc    *KRPC
	outChan chan string
}

type KRPC struct {
	dhtNode *KNode
	tid     uint32
}

type KRPCMSG struct {
	T   string
	Y   string
	Ags interface{}
	add *net.UDPAddr
}

type Query struct {
	Q string
	A map[string]interface{}
}

type Response struct {
	R map[string]interface{}
}


func (dhtNode *KNode) Run() {

}
func main() {

}
