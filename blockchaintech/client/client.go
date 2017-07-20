package main

import (
    "github.com/ethereum/go-ethereum/ethclient"
    "fmt"
)

func main() {
    c, e := ethclient.Dial("http://localhost:8545")

    if e != nil {
        fmt.Println("Error ", e)
        return
    }
    c.BlockByHash(nil, "")
}
