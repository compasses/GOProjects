package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/juju/errors"
	"github.com/zeebo/bencode"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"runtime"
	"time"
)

// refer :https://github.com/jyootai/DhtBT/

var BOOTSTRAP []string = []string {
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
	conn    *net.UDPConn
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
	T    string
	Y    string
	Args interface{}
	addr *net.UDPAddr
}

type Query struct {
	Q string
	A map[string]interface{}
}

type Response struct {
	R map[string]interface{}
}

func (dhtNode *KNode) Run() {
	dhtNode.log.Println(fmt.Sprintf("Current Node ID is %s ", dhtNode.node.id))
	dhtNode.log.Println(fmt.Sprintf("DhtBT %s is runing...", dhtNode.network.conn.LocalAddr().String()))

}

func (network *Network) GetInfohash() {
	b := make([]byte, 1000)
	for {
		_, addr, err := network.conn.ReadFromUDP(b)
		if err != nil {
			continue
		}
		network.dhtNode.krpc.DecodePackage(string(b), addr)
	}
}

func (krpc *KRPC) DecodePackage(data string, addr *net.UDPAddr) error {
	val := make(map[string]interface{})
	if err := bencode.DecodeString(data, &val); err != nil {
		return err
	} else {
		var ok bool
		msg := new(KRPCMSG)
		msg.T, ok = val['t'].(string)
		if !ok {
			err = errors.New("Do not have transaction ID")
			return err
		}
		msg.Y, ok = val['y'].(string)
		if !ok {
			err = errors.New("Do know message type ")
			return err
		}

		msg.addr = addr
		switch msg.Y {
		case "q":
			query := new(Query)
			query.Q = val["q"].(string)
			query.A = val["a"].(map[string]interface{})
			msg.Args = query
			krpc.Query(msg)
		case "r":
			res := new(Response)
			res.R = val["r"].(map[string]interface{})
			msg.Args = res
			krpc.Response(msg)
		}
	}
	return nil
}

func (krpc *KRPC) Response(msg *KRPCMSG) {
	if res, ok := msg.Args.(*Response); ok {
		if nodestr, ok := res.R["nodes"].(string); ok {
			nodes := ParseBytesStream([]byte(nodestr))
			for _, v := range nodes {
				krpc.dhtNode.routing.InsertNode(v)
			}
		}
	}
}

func (routing *Routing) InsertNode(other *NodeInfo) {
	if routing.isSelf(other) {
		return
	}

	if routing.table[1].Len() < 8 {
		routing.table[1].Add(other)
	}

	routing.table[0].Add(other)
}

func (table *Bucket) Add(n *NodeInfo) {
	table.nodes = append(table.nodes, n)
	table.Updatetime(n)
}

func (bucket *Bucket) Updatetime(n *NodeInfo ) {
	bucket.lastChange = time.Now()
	n.lastSeen = time.Now()
}

func (id *Id) HexString() string {
	return fmt.Sprintf("%x", id)
}

func (routing *Routing) isSelf(other *NodeInfo) bool {
	return (routing.selfNode.node.id.CompareTo(other.id) == 0)
}

func (id Id) CompareTo(other Id) int {
	s1 := id.HexString()
	s2 := other.HexString()
	if s1 > s2 {
		return 1
	} else if s1 == s2 {
		return 0
	} else {
		return -1
	}
}


func ParseBytesStream(data []byte) []*NodeInfo {
	var nodes []*NodeInfo = nil
	for j := 0; j < len(data); j = j + 26 {
		if j+26 > len(data) {
			break
		}
		kn := data[j : j+26]
		node := new(NodeInfo)
		node.id = Id(kn[0:20])
		node.ip = kn[20:24]
		port := kn[24:26]
		node.port = int(port[0])<<8 + int(port[1])
		nodes = append(nodes, node)
	}
	return nodes
}

func (krpc *KRPC) Query(msg *KRPCMSG) {
	if query, ok := msg.Args.(*Query); ok {
		queryNode := new(NodeInfo)

	}
}

func (krpc *KRPC) EncodingNodeResult(tid string, token string, nodes []byte) (string, error) {
	v := make(map[string]interface{})
	v['t'] = tid
	v['y'] = "r"

	args := make(map[string]string)
	args["id"] = krpc.dhtNode.node.id.String()
	if token != "" {
		args["token"] = token
	}
	args["nodes"] = bytes.NewBuffer(nodes).String()
	v["r"] = args
	s, err := bencode.EncodeString(v)
	return s, err
}

func (id Id) String() string {
	return hex.EncodeToString(id)
}

var nums int = 0

func OutHash(master chan string) {
	for {
		select {
		case infohash := <-master:
			fmt.Println("get hash:", infohash)
			nums++
			fmt.Println("total nums:", nums)
		}
	}
}

func main() {
	cpu := runtime.NumCPU()
	fmt.Println("Number CPU is ", cpu)
	runtime.GOMAXPROCS(cpu)
	master := make(chan string)

	for i := 0; i < cpu; i++ {
		go Executing(cpu, master)
	}

	OutHash(master)
}

func NewDHTNode(master chan string, logger *io.Writer) *KNode {
	dhtNode := new(KNode)
	dhtNode.log = log.New(logger, "DHT_TEST", log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	dhtNode.node = NewNodeInfo()
	dhtNode.routing = NewRouting(dhtNode)
	dhtNode.network = NewNetwork(dhtNode)
	dhtNode.krpc = NewKrpc(dhtNode)
	dhtNode.outChan = master
	return dhtNode
}

func GenerateId() Id {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	h := sha1.New()
	io.WriteString(h, time.Now().String())
	io.WriteString(h, string(random.Int()))
	return h.Sum(nil)
}

func NewNodeInfo() *NodeInfo {
	node := new(NodeInfo)
	id := GenerateId()
	node.id = id
	return node
}
func NewBucket() *Bucket {
	b := new(Bucket)
	b.nodes = nil
	b.lastChange = time.Now()
	return b
}

func NewBucket2() *Bucket {
	b := new(Bucket)
	b.nodes = nil
	b.lastChange = time.Now()
	return b
}
func NewRouting(dhtNode *KNode) *Routing {
	routing := new(Routing)
	routing.selfNode = dhtNode
	bucket1 := NewBucket()
	bucket2 := NewBucket2()
	routing.table = make([]*Bucket, 2)
	routing.table[0] = bucket1
	routing.table[1] = bucket2
	return routing
}

func (bucket *Bucket) Len() int {
	return len(bucket.nodes)
}

func NewNetwork(dhtNode *KNode) *Network {
	network := new(Network)
	network.dhtNode = dhtNode
	addr := new(net.UDPAddr)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	network.conn = conn

	laddr := conn.LocalAddr().(*net.UDPAddr)
	network.dhtNode.node.ip = laddr.IP
	network.dhtNode.node.port = laddr.Port
	return network
}

func NewKrpc(dhtNode *KNode) *KRPC {
	krpc := new(KRPC)
	krpc.dhtNode = dhtNode
	return krpc
}

func Executing(cpu int, master chan string) {
	dhtNode := NewDHTNode(master, os.Stdout)
	dhtNode.node = NewNodeInfo()

}
