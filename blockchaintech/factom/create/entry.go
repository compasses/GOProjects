package main

import (
    "log"
    "time"

    "github.com/FactomProject/factom"
)

func createNewChain() {
    e := factom.Entry{}
    e.ExtIDs = append(e.ExtIDs, []byte("MyChain"), []byte("12345"))
    e.Content = []byte("Hello Factom!")

    c := factom.NewChain(&e)
    log.Println("Creating new Chain:", c.ChainID)

    var str string
    factom.RpcConfig.FactomdServer = "localhost:8088"

    ec, err := factom.GetECAddress("Es3h3RWmtmrE2XVUHU39d7DGPGRjbfeanNaZJebzxc6ZyDch3hFt")
    if err != nil {
        log.Println("Got error ", err)
        return
    }

    log.Println("got address ", ec)

    if str, err := factom.CommitChain(c, ec); err != nil {
        log.Fatal(err, str)
    }
    log.Println("Got ", str)

    time.Sleep(10 * time.Second)
    if str, err := factom.RevealChain(c); err != nil {
        log.Fatal(err, str)
    }
}

func createEntry() {
    e := factom.Entry{}
    e.ChainID = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    e.ExtIDs = append(e.ExtIDs, []byte("hello"))
    e.Content = []byte("Hello Factom!")

    ec, err := factom.GetECAddress("Es3h3RWmtmrE2XVUHU39d7DGPGRjbfeanNaZJebzxc6ZyDch3hFt")
    if err != nil {
        log.Println("Got error ", err)
        return
    }

    log.Println("got address ", ec)
    var str string
    factom.RpcConfig.FactomdServer = "localhost:8088"

    if str, err := factom.CommitEntry(&e, ec); err != nil {
        log.Fatal(err)
        log.Println(str)
    }

    log.Println("str" , str)

    time.Sleep(10 * time.Second)
    if str, err := factom.RevealEntry(&e); err != nil {
        log.Fatal(err)
        log.Println(str)
    }
    log.Println("str" , str)
}

func main() {
   createNewChain()
}
