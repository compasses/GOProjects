package main

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
    "log"
    "os"
)

func MyProtocol() p2p.Protocol {
    return p2p.Protocol{ // 1.
        Name:    "MyProtocol",                                                    // 2.
        Version: 1,                                                               // 3.
        Length:  1,                                                               // 4.
        Run:     func(peer *p2p.Peer, ws p2p.MsgReadWriter) error { return nil }, // 5.
    }
}

const messageId = 0

type Message string


func msgHandler(peer *p2p.Peer, ws p2p.MsgReadWriter) error {
    for {
        msg, err := ws.ReadMsg()
        if err != nil {
            return err
        }

        var myMessage Message
        err = msg.Decode(&myMessage)
        if err != nil {
            // handle decode error
            continue
        }

        switch myMessage {
        case "foo":
            err := p2p.SendItems(ws, messageId, "bar")
            if err != nil {
                return err
            }
        default:
            log.Println("recv:", myMessage)
        }
    }

    return nil
}

func main() {
	nodekey, _ := crypto.GenerateKey()
    conf := p2p.Config{
        MaxPeers:   10,
        PrivateKey: nodekey,
        Name:       "my node name",
        ListenAddr: ":30300",
        Protocols:  []p2p.Protocol{},
    }

    log.Println("Using private key: ", nodekey)

	srv := p2p.Server{
		Config: conf,
	}

    if err := srv.Start(); err != nil {
        log.Println(err)
        os.Exit(1)
    }

    select {}
}
