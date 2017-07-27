package main

import (
    "os/user"
    "os"
    log "github.com/sirupsen/logrus"

    "github.com/FactomProject/factom/wallet"
    "github.com/FactomProject/factom"
    "net/http"
    "fmt"
    "bytes"
    "encoding/json"
    "io/ioutil"
    "github.com/FactomProject/factomd/common/primitives/random"
    "time"
    "github.com/FactomProject/factomd/anchor"
    "github.com/FactomProject/factomd/common/primitives"
    "golang.org/x/crypto/ed25519"
)
var contextLogger *log.Entry

func init() {
    log.SetLevel(log.InfoLevel)
    log.WithFields(log.Fields{"package": "state"})
    contextLogger = log.WithFields(log.Fields{
        "common": "entry create test",
    })

}

func getHomeDir() string{
    usr, err := user.Current()
    if err == nil {
        return usr.HomeDir
    } else {
    }
    return os.Getenv("HOME")

}

func openWallet(path string) *wallet.Wallet {
    w, err := wallet.NewOrOpenBoltDBWallet(path)
    if err != nil {
        contextLogger.Fatal(err)
    }

    // open and add a transaction database to the wallet object.
    txdb, err := wallet.NewTXBoltDB(fmt.Sprint(getHomeDir(), "/.factom/wallet/factoid_blocks.cache"))
    if err != nil {
        contextLogger.Println("Could not add transaction database to wallet:", err)
    } else {
        w.AddTXDB(txdb)
    }

    return w
}

func newEntryReq(ecpub, chainId, hash, content string) string  {
    hashs := []string{hash}
    body := map[string]interface{}{
            "entry": map[string]interface{}{
                "chainid": chainId,
                "extids":hashs,
                "content": content,
            },
            "ecpub": ecpub,
    }

    req := factom.NewJSON2Request("compose-entry", 0, body)

    str, err := factom.EncodeJSONString(req)

    if err != nil {
        contextLogger.Info("Error on creating json ", err)
        os.Exit(1)
    }

    return str
}

func newEntry(chainid string, hash []byte, content string) *factom.Entry{
    entry := new(factom.Entry)
    entry.ChainID = chainid
    entry.Content = []byte(content)
    bta := make([][]byte, 0)
    bta = append(bta, []byte(hash))
    entry.ExtIDs = bta

    return entry
}

func getNewEntry() *factom.Entry {
    chainid := "40fa997b215025e1a24d9a5a3c2ab9526589ce369e0bc4859e37023312822e5b"
    hash := []byte("40fa997b215025e1a24d9a5a3c2ab9526589ce369e0bc4859e37023312822e5b" + random.RandomString())
    content := "random test, 40fa997b215025e1a24d9a5a3c2ab9526589ce369e0bc4859e37023312822e5b, random test" + random.RandomString()
    return newEntry(chainid, hash, content)
}

var anchorRec = `
{
	"AnchorRecordVer": 1,
	"DBHeight": 8,
	"KeyMR": "637b6010cb6121f76c65b200a6cf94cb6655881fb4cac48979f8950e7a349da1",
	"RecordHeight": 9,
	"Bitcoin": {
		"Address": "1K2SXgApmo9uZoyahvsbSanpVWbzZWVVMF",
		"TXID": "b73b38b8af43f4dbaeb061f158d4bf5004b40216b30acd3beca43fae1ba6d1b7",
		"BlockHeight": 372579,
		"BlockHash": "00000000000000000589540fdaacf4f6ba37513aedc1033e68a649ffde0573ad",
		"Offset": 1185
	}}`

var private = "e89f1216745b61d056b5297be7bdecd7b82966feb2ae482f686e127bd1b2ff80"
var public = "672f1517608fe33fd6be113412afd22b8adc8ece0f4b823071f2810dafe4f153"

func getNewAnchorRecord(dbheight uint32, keyMR string) *anchor.AnchorRecord  {

    ar, err := anchor.UnmarshalAnchorRecord([]byte(anchorRec))
    if err != nil {
        contextLogger.Error("error on unmarshal record data", err)
        os.Exit(1)
    }
    ar.DBHeight = dbheight
    ar.RecordHeight = dbheight + 1
    ar.KeyMR = keyMR

    contextLogger.Info("New record is ", ar)
    return ar
}

func signRecord(ar *anchor.AnchorRecord) ([]byte, []byte) {
    //prv := primitives.RandomPrivateKey()
    prv , err := primitives.NewPrivateKeyFromHex(private)

    if err != nil {
        fmt.Println("New private key error", err)
        os.Exit(1)
    }

    fmt.Println("private string ", prv.PrivateKeyString())
    fmt.Println("Public string ", prv.PublicKeyString())

    raw, sign, err := ar.MarshalAndSignV2(prv)

    if err != nil {
        fmt.Println("Error on sign is ", err)
        os.Exit(1)
    }

    fmt.Println("got raw is ", raw)
    fmt.Println("got sign is ", sign)
    sig := prv.Sign(raw)
    fmt.Println("got sign is ", sig.GetSignature())
    fmt.Println("verfiy msg result is ", prv.Pub.Verify(raw, sig.GetSignature()))
    return raw, sign
}


func getNewEntryForAnchor() *factom.Entry {
    chainid := "df3ade9eec4b08d5379cc64270c30ea7315d8a8a1a69efe2b98a60ecdd69e604"

    ar := getNewAnchorRecord(48, "0fbe914f3819cbd8f9bfc50be4dd29dc807a7768f135b62839667dbfff1dc0fa")

    raw, sign  := signRecord(ar)

    return newEntry(chainid, sign, string(raw))
}


func startPushEntryCommit(wallet *wallet.Wallet, ecpub string)  {
    contextLogger.Info("Start commit entry ...")
    ecaddr, err := wallet.GetECAddress(ecpub)
    if err != nil {
        contextLogger.Error("open wallet error ", err)
        os.Exit(1)
    }

//    en := getNewEntry()

    en := getNewEntryForAnchor()

    commit, err := factom.ComposeEntryCommit(en, ecaddr)

    if err != nil {
        contextLogger.Error("Fatal error on commit entry ", err)
    }
    commitBody, err := factom.EncodeJSON(commit)

    if err != nil {
        contextLogger.Error("Encode error ", commitBody)
    }

    httpClient := http.DefaultClient

    contextLogger.Info("do commit ", string(commitBody))
    re, err := http.NewRequest("POST", fmt.Sprintf("http://%s/v2", "localhost:8088"), bytes.NewBuffer(commitBody))

    if err != nil {
        contextLogger.Error("error happened, for entry commit ", err)
        os.Exit(1)
    }

    resp, err := httpClient.Do(re)
    if err != nil {
        contextLogger.Error("Error for http request ", err)
        os.Exit(1)
    }

    if resp.StatusCode == http.StatusUnauthorized {
        contextLogger.Error("Factomd username/password incorrect.  Edit factomd.conf or\ncall factom-cli with -factomduser=<user> -factomdpassword=<pass>")
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)

    r := factom.NewJSON2Response()
    if err := json.Unmarshal(body, r); err != nil {
        contextLogger.Error("Error on http request parse body", err)
    }

    contextLogger.Info("Got response ", r)
    time.Sleep(2000)
    rev, err := factom.ComposeEntryReveal(en)
    if err != nil {
        contextLogger.Info("Got error on entry compose ", err)
    }

    revBody, err := factom.EncodeJSON(rev)

    contextLogger.Println("Do reveal ", string(revBody))
    if err != nil {
        contextLogger.Error("Encode error ", revBody)
    }

    re, err = http.NewRequest("POST", fmt.Sprintf("http://%s/v2", "localhost:8088"), bytes.NewBuffer(revBody))

    if err != nil {
        contextLogger.Error("error happened, for entry revl ", err)
        os.Exit(1)
    }

    resp, err = httpClient.Do(re)
    if err != nil {
        contextLogger.Error("Error for http request ", err)
        os.Exit(1)
    }

    if resp.StatusCode == http.StatusUnauthorized {
        log.Error("Factomd username/password incorrect.  Edit factomd.conf or\ncall factom-cli with -factomduser=<user> -factomdpassword=<pass>")
    }
    defer resp.Body.Close()

    body, err = ioutil.ReadAll(resp.Body)

    r = factom.NewJSON2Response()
    if err := json.Unmarshal(body, r); err != nil {
        contextLogger.Error("Error on http request parse body", err)
    }

    contextLogger.Info("Got response ", r)

}




func main() {
    //open wallet
    wallet := openWallet("/Users/jet/.factom/wallet/factom_wallet.db")
    contextLogger.Info("got wallet %q", wallet)
    ecpub := "EC1nEmE4oA7T5KzicS6P7Gzoyi2d73rxjtm3VZCvBVZgDwKE4pLz"
    ecs, err := wallet.GetAllECAddresses()
    if err != nil {
        contextLogger.Info("got error get all ec address, ", err)
    }
    contextLogger.Info("addresses ", ecs)
    fmt.Println("use random bytes is ", random.RandByteSliceOfLen(ed25519.PrivateKeySize))
    //getNewEntryForAnchor()
    startPushEntryCommit(wallet, ecpub)
    //go func() {
    //    for {
    //        time.Sleep(20000)
    //        startPushEntryCommit(wallet, ecpub)
    //    }
    //}()
    select {

    }

    // new entryRequest
}
