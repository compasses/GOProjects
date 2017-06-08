package main

//a file server, can send the file which assigned on command
import (
	"syscall"
	"os"
	"fmt"
	"net"
	"flag"

)

var filename = flag.String("f", ".", "file to send")
var listenerAddr = flag.String("p", ":9192", "listen ports")


func main() {
	flag.Parse()
	listenAddr, err := net.ResolveTCPAddr("tcp", *listenerAddr)
	checkPanic(err)
	listener, err := net.ListenTCP("tcp", listenAddr)
	checkPanic(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection '%s'\n", err)
			continue
		}		
		go handleConn(conn)
	}
	
}

func handleConn(c net.Conn) {
	fmt.Println("Got connection: ", c.LocalAddr(), "remote ", c.RemoteAddr())
    defer c.Close()
	buf := make([]byte, 1024)
    for {		
        // read from the connection
        n, err := c.Read(buf)
		checkIgnore(err)
		if n <= 0 {
			fmt.Println("connection break...")
			return
		}

		if n > 0 {
			//should be a file name
			f, err := os.Open(string(buf[:n]))
			checkIgnore(err)
			if f != nil {
				//send this file to the client
				sendFile(f, c)
				return
			} else {
				c.Write([]byte("error"));
				fmt.Println("Open file error ", err)
				return
			}
		}
    }
}

func sendFile(file *os.File, c net.Conn) {
	fStat, err := file.Stat()
	checkIgnore(err)
	fmt.Println("Start to send file, name ", file.Name(), "size ", fStat.Size())

	addr, err := syscall.Mmap(int(file.Fd()), 0,  int(fStat.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println("Send file error when mmap ", err)
		return
	}
	n, err := c.Write(addr)
	if err != nil {
		fmt.Println("write error , err")
	}
	fmt.Println("write finish :", n)
}

func checkPanic(err error) {
	if err != nil {
		fmt.Println("got error:", err)
		os.Exit(1)
	}
}

func checkIgnore(err error) {
	if err != nil {
		fmt.Println("got error:", err)
	}
}
