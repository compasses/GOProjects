package main

import (
    "fmt"
    "github.com/btcsuite/btcutil"
    "github.com/btcsuite/btcd/chaincfg"
    "github.com/btcsuite/btcd/txscript"
)

func main() {
    // Parse the address to send the coins to into a btcutil.Address
    // which is useful to ensure the accuracy of the address and determine
    // the address type.  It is also required for the upcoming call to
    // PayToAddrScript.
    addressStr := "myNyL1X25ikv3BUqhT1p1TYEdZLnw1JNiG"
    address, err := btcutil.DecodeAddress(addressStr, &chaincfg.TestNet3Params)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Create a public key script that pays to the address.
    script, err := txscript.PayToAddrScript(address)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("Script Hex: %x\n", script)

    disasm, err := txscript.DisasmString(script)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Script Disassembly:", disasm)
}
