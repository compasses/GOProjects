package main

import (
	"unsafe"
	"fmt"
	"os"
)

// DO NOT EDIT. Copied from the "bolt" package.
const bucketLeafFlag = 0x01

// DO NOT EDIT. Copied from the "bolt" package.
type pgid uint64

// DO NOT EDIT. Copied from the "bolt" package.
type txid uint64


// DO NOT EDIT. Copied from the "bolt" package.
type meta struct {
	magic    uint32
	version  uint32
	pageSize uint32
	flags    uint32
	root     bucket
	freelist pgid
	pgid     pgid
	txid     txid
	checksum uint64
}

// DO NOT EDIT. Copied from the "bolt" package.
type bucket struct {
	root     pgid
	sequence uint64
}

// DO NOT EDIT. Copied from the "bolt" package.
type page struct {
	id       pgid
	flags    uint16
	count    uint16
	overflow uint32
	ptr      uintptr
}

func pageInBuffer(b []byte, id pgid) *page {
	return (*page)(unsafe.Pointer(&b[id*pgid(4096)]))
}

func (p *page) meta() *meta {
	return (*meta)(unsafe.Pointer(&p.ptr))
}

func showPageInfo(b []byte, id pgid) {
	p := pageInBuffer(b, id)
	fmt.Printf("page id %d, flags %d, count %d, overflow %d, ptr %d\n", p.id, p.flags, p.count, p.overflow, p.ptr)
	fmt.Printf("Meta info \n")
	m := p.meta()
	
	w := os.Stdout
	fmt.Fprintf(w, "Version:    %d\n", m.version)
	fmt.Fprintf(w, "Page Size:  %d bytes\n", m.pageSize)
	fmt.Fprintf(w, "Flags:      %08x\n", m.flags)
	fmt.Fprintf(w, "Root:       <pgid=%d>\n", m.root.root)
	fmt.Fprintf(w, "Freelist:   <pgid=%d>\n", m.freelist)
	fmt.Fprintf(w, "HWM:        <pgid=%d>\n", m.pgid)
	fmt.Fprintf(w, "Txn ID:     %d\n", m.txid)
	fmt.Fprintf(w, "Checksum:   %016x\n", m.checksum)
	fmt.Fprintf(w, "\n")

}

func main() {
	
	db, err := os.Open("./newDB");
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	finfo, _ := db.Stat()
	sz := finfo.Size()
	
	buf := make([]byte, sz)
	
	n, err := db.ReadAt(buf[:], 0)
	if n < 0x1000 || err != nil {
		fmt.Println("Read error, got len  ", n, " error ", err)
	}
	fmt.Println("Read error, got len  ", n, " error ", err)
	showPageInfo(buf[:], 7)
	
}