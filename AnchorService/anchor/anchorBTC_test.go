package anchor_test

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"testing"
)

func TestPayScript(t *testing.T) {
	addressStr := "myNyL1X25ikv3BUqhT1p1TYEdZLnw1JNiG"
	address, err := btcutil.DecodeAddress(addressStr, &chaincfg.TestNet3Params)
	if err != nil {
		fmt.Println("decodeaddress", err)
		return
	}

	fmt.Printf("Address %q\n", address.EncodeAddress())

	script, err := txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
		AddData(address.ScriptAddress()).AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).
		Script()

	if err != nil {
		fmt.Println("paytoaddrerror", err)
		fmt.Println("returen is ", script)
		return
	}

	switch address := address.(type) {
	case *btcutil.AddressScriptHash:
		break
	case *btcutil.AddressPubKeyHash:
		fmt.Println("yes it's the addre", address)
		if address == nil {
			fmt.Println("Let's breakdown ,,,")
		}
	case *btcutil.AddressPubKey:

	}

	//Create a public key script that pays to the address.
	//script, err := txscript.PayToAddrScript(address)
	//if err != nil {
	//    fmt.Println("paytoaddrerror", err)
	//    fmt.Println("returen is ", script)
	//    return
	//}
	fmt.Printf("Script Hex: %x\n", script)

	disasm, err := txscript.DisasmString(script)
	if err != nil {
		fmt.Println("disasmstring", err)
		return
	}
	fmt.Println("Script Disassembly:", disasm)
}
