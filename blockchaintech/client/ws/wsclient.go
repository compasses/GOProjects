package main

import (
    "time"
    "fmt"
    "github.com/ethereum/go-ethereum/rpc"
    "math/big"
    "context"
)


type Block struct {
    Number *big.Int
}

func main() {
    client, e := rpc.Dial("ws://127.0.0.1:8546")
    subch := make(chan Block)

    if e != nil {
        fmt.Println("Error ", e)
        return
    }
    // Ensure that subch receives the latest block.
    go func() {
        for i := 0; ; i++ {
            if i > 0 {
                time.Sleep(2 * time.Second)
            }
            subscribeBlocks(client, subch)
        }
    }()

    // Print events from the subscription as they arrive.
    for block := range subch {
        fmt.Println("latest block:", block.Number)
    }
}


// subscribeBlocks runs in its own goroutine and maintains
// a subscription for new blocks.
func subscribeBlocks(client *rpc.Client, subch chan Block) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Subscribe to new blocks.
    sub, err := client.EthSubscribe(ctx, subch, "newBlockFilter")
    if err != nil {
        fmt.Println("subscribe error:", err)
        //return
    }

    // The connection is established now.
    // Update the channel with the current block.
    var lastBlock Block
    if err := client.CallContext(ctx, &lastBlock, "eth_getBlockByNumber", "latest"); err != nil {
        fmt.Println("can't get latest block:", err)
        return
    }
    subch <- lastBlock

    // The subscription will deliver events to the channel. Wait for the
    // subscription to end for any reason, then loop around to re-establish
    // the connection.
    fmt.Println("connection lost: ", <-sub.Err())
}
