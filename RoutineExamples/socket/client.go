package main

//get data from server, and save it a given filename.
import (
	"time"
	"os"
	"log"
	"net"
	"flag"
)

var saveFile = flag.String("f", "receive.file", "file name to save")
var serverAddr = flag.String("s", "localhost:9192", "server address")

func main() {
	flag.Parse()
	
	log.Println("begin dial...")
    conn, err := net.Dial("tcp", *serverAddr)
    if err != nil {
        log.Println("dial error:", err)
        return
    }

    var buf = make([]byte, 1024*1024)
	
	log.Println("Got connect")
	conn.Write([]byte(*saveFile))	

	file , err := os.OpenFile(*saveFile, os.O_RDWR|os.O_CREATE, 0666)
	
	if err != nil {
		log.Println("Create file error ", err)
		return
	}
	//conn.SetReadDeadline(time.Now().Add(time.Microsecond * 100))
	var totalLen int
	now := time.Now()
	defer func() {
		log.Println("Time used ", time.Since(now))
		log.Println("Got total bytes ", totalLen)
		file.Close()
		conn.Close()
	}()
	
	for {
		n, err := conn.Read(buf)
	    if err != nil {
	        log.Println("read error:", err)
			return
	    }
		if string(buf[:n]) == "error" {
			log.Println("Server not found this file...")
			return
		}

		if n <= 0 {
			log.Println("got len ", n)
			return
		}
		totalLen += n
		file.Write(buf[:n])	
		//log.Println("Got total bytes ", totalLen)
	}
	
	
}


