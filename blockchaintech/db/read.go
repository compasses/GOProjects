package main

import (
    "github.com/syndtr/goleveldb/leveldb"
    "log"
)


var dbpath string = "/Users/jet/Library/Ethereum/geth/chaindata"


func readBlock() {
    db, err := leveldb.OpenFile(dbpath, nil)
    if err != nil {
        log.Println("Got error ", err)
    }

    iter := db.NewIterator(nil, nil)
    i := 0
    for iter.Next() {
        // Remember that the contents of the returned slice should not be modified, and
        // only valid until the next call to Next.
        iter.Key()
        iter.Value()
        i ++
       // log.Printf("Key %s\n", key)
        //log.Printf("Value %s\n", value)
    }
    iter.Release()
    err = iter.Error()

    if err != nil {
        log.Print("GO error \n", err)
    }
    log.Println("total count ", i)
    defer db.Close()
}
