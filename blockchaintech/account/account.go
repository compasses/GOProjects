package main

import (
    "github.com/ethereum/go-ethereum/accounts"
    "github.com/ethereum/go-ethereum/accounts/keystore"
    "io/ioutil"
    "fmt"
    "github.com/ethereum/go-ethereum/common"
)

func AccountCreateTest() {
    scryptN := keystore.StandardScryptN
    scryptP := keystore.StandardScryptP
    keydir, err := ioutil.TempDir("", "go-ethereum-keystore")
    if err != nil {
        fmt.Println("errors ", err)
        return
    }
    backends := []accounts.Backend{
        keystore.NewKeyStore(keydir, scryptN, scryptP),
    }

    am := accounts.NewManager(backends...);
    fmt.Println("got manager is ", am)

    keystore := keystore.NewKeyStore("./", scryptN, scryptP)


    // Create a new account with the specified encryption passphrase.
    newAcc, _ := keystore.NewAccount("Creation password");

    // Export the newly created account with a different passphrase. The returned
    // data from this method invocation is a JSON encoded, encrypted key-file.
    jsonAcc, _ := keystore.Export(newAcc, "Creation password", "Export password")
    fmt.Printf("json account %q\n", jsonAcc)

    // Update the passphrase on the account created above inside the local keystore.
    keystore.Update(newAcc, "Creation password", "Update password");

    // Delete the account updated above from the local keystore.
    keystore.Delete(newAcc, "Update password");

    // Import back the account we've exported (and then deleted) above with yet
    // again a fresh passphrase.
    impAcc, _ := keystore.Import(jsonAcc, "Export password", "Import password");
    fmt.Printf("use the import accout %q\n", impAcc)
}

func TxSignTest() {
    scryptN := keystore.StandardScryptN
    scryptP := keystore.StandardScryptP
    ks :=  keystore.NewKeyStore("./", scryptN, scryptP)
    backends := []accounts.Backend{
       ks,
    }

    am := accounts.NewManager(backends...);
    fmt.Println("got manager is ", am)
    txHash := common.HexToHash("0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
    fmt.Println("Got hash ", txHash)
    signer, err := ks.NewAccount("Signed User!")
    if err != nil {
        fmt.Println("Create account  error ", err)
    }
    sg, err := ks.SignHashWithPassphrase(signer, "Signed User!", txHash.Bytes())
    if err != nil {
        fmt.Println("Signed error ", err)
    }
    fmt.Println("Got signature ", sg)

    ks.Unlock(signer, "Signed User!")
    sg, err = ks.SignHash(signer, txHash.Bytes())
    if err != nil {
        fmt.Println("Signed error ", err)
    }
    fmt.Println("Got signature ", sg)

    //signature, _ := ks.SignHash()
}


func main() {
    TxSignTest()
}
