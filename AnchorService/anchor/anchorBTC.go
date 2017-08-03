package anchor

import (
	"AnchorService/common"
	"AnchorService/util"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/FactomProject/factomd/anchor"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcrpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

type balance struct {
	unspentResult btcjson.ListUnspentResult
	address       btcutil.Address
	wif           *btcutil.WIF
}

type AnchorBTC struct {
	service             *AnchorService
	toAnchorInfo        map[string]*anchor.AnchorRecord
	balances            []balance // unspent balance & address & its WIF
	dclient, wclient    *btcrpcclient.Client
	cfg                 *util.AnchorServiceCfg
	fee                 btcutil.Amount // tx fee for written into btc
	defaultAddress      btcutil.Address
	walletLocked        bool
	confirmationsNeeded int
}

func NewAnchorBTC() *AnchorBTC {
	btc := &AnchorBTC{
		toAnchorInfo: make(map[string]*anchor.AnchorRecord),
	}
	return btc
}

func (anchorBTC *AnchorBTC) PlaceAnchor(msg common.DirectoryBlockAnchorInfo) {
	anchorRec := new(anchor.AnchorRecord)
	btc := new(anchor.BitcoinStruct)

	anchorRec.Bitcoin = btc
	anchorRec.KeyMR = msg.KeyMR.String()
	anchorRec.DBHeight = msg.DBHeight
	anchorRec.AnchorRecordVer = 1
	anchorBTC.doTransaction(anchorRec, msg.KeyMR)
}

func (anchorBTC *AnchorBTC) InitRPCClient() error {
	log.Debug("init anchorBTC RPC client")
	cfg := util.ReadConfig()
	certHomePath := cfg.Btc.CertHomePath
	rpcClientHost := cfg.Btc.RpcClientHost
	rpcClientEndpoint := cfg.Btc.RpcClientEndpoint
	rpcClientUser := cfg.Btc.RpcClientUser
	rpcClientPass := cfg.Btc.RpcClientPass
	certHomePathBtcd := cfg.Btc.CertHomePathBtcd
	rpcBtcdHost := cfg.Btc.RpcBtcdHost
	anchorBTC.cfg = cfg

	// Connect to local btcwallet RPC server using websockets.
	ntfnHandlers := anchorBTC.createBtcwalletNotificationHandlers()
	certHomeDir := btcutil.AppDataDir(certHomePath, false)
	log.Debug("btcwallet.cert.home=", certHomeDir)
	certs, err := ioutil.ReadFile(filepath.Join(certHomeDir, "rpc.cert"))
	if err != nil {
		return fmt.Errorf("cannot read rpc.cert file: %s\n", err)
	}
	connCfg := &btcrpcclient.ConnConfig{
		Host:         rpcClientHost,
		Endpoint:     rpcClientEndpoint,
		User:         rpcClientUser,
		Pass:         rpcClientPass,
		Certificates: certs,
	}

	wclient, err := btcrpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		return fmt.Errorf("cannot create rpc client for btcwallet: %s\n", err)
	}
	anchorBTC.wclient = wclient

	log.Debug("successfully created rpc client for btcwallet")

	// Connect to local btcd RPC server using websockets.
	dntfnHandlers := anchorBTC.createBtcdNotificationHandlers()
	certHomeDir = btcutil.AppDataDir(certHomePathBtcd, false)
	log.Debug("btcd.cert.home=", certHomeDir)
	certs, err = ioutil.ReadFile(filepath.Join(certHomeDir, "rpc.cert"))
	if err != nil {
		return fmt.Errorf("cannot read rpc.cert file for btcd rpc server: %s\n", err)
	}
	dconnCfg := &btcrpcclient.ConnConfig{
		Host:         rpcBtcdHost,
		Endpoint:     rpcClientEndpoint,
		User:         rpcClientUser,
		Pass:         rpcClientPass,
		Certificates: certs,
	}

	dclient, err := btcrpcclient.New(dconnCfg, &dntfnHandlers)
	if err != nil {
		return fmt.Errorf("cannot create rpc client for btcd: %s\n", err)
	}
	anchorBTC.dclient = dclient
	log.Debug("successfully created rpc client for btcd")

	err = anchorBTC.initWallet()
	if err != nil {
		log.Fatal("Init Wallet error ", err)
	}

	return nil
}

func (anchorBTC *AnchorBTC) unlockWallet(timeoutSecs int64) error {
	err := anchorBTC.wclient.WalletPassphrase(anchorBTC.cfg.Btc.WalletPassphrase, int64(timeoutSecs))
	if err != nil {
		return fmt.Errorf("cannot unlock wallet with passphrase: %s", err)
	}
	anchorBTC.walletLocked = false
	return nil
}

func (anchorBTC *AnchorBTC) initWallet() error {
	anchorBTC.fee, _ = btcutil.NewAmount(anchorBTC.cfg.Btc.BtcTransFee)
	anchorBTC.walletLocked = true
	err := anchorBTC.updateUTXO()
	if err == nil && len(anchorBTC.balances) > 0 {
		anchorBTC.defaultAddress = anchorBTC.balances[0].address
	}
	return err
}

func (anchorBTC *AnchorBTC) updateUTXO() error {
	log.Info("updateUTXO: walletLocked=", anchorBTC.walletLocked)
	anchorBTC.balances = make([]balance, 0, 200)

	err := anchorBTC.unlockWallet(int64(6)) //600
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	unspentResults, err := anchorBTC.wclient.ListUnspentMin(anchorBTC.cfg.Anchor.ConfirmationsNeeded) //minConf=1
	if err != nil {
		return fmt.Errorf("cannot list unspent. %s", err)
	}
	log.Info("updateUTXO: unspentResults.len=", len(unspentResults))
	log.Debug("unspent result is ", unspentResults)

	if len(unspentResults) > 0 {
		var i int
		for _, b := range unspentResults {
			if b.Amount > anchorBTC.fee.ToBTC() {
				anchorBTC.balances = append(anchorBTC.balances, balance{unspentResult: b})
				i++
			}
		}
	}
	log.Info("updateUTXO: balances.len=", len(anchorBTC.balances))

	for i, b := range anchorBTC.balances {
		addr, err := btcutil.DecodeAddress(b.unspentResult.Address, &chaincfg.TestNet3Params)
		if err != nil {
			return fmt.Errorf("cannot decode address: %s", err)
		}
		anchorBTC.balances[i].address = addr

		wif, err := anchorBTC.wclient.DumpPrivKey(addr)
		if err != nil {
			return fmt.Errorf("cannot get WIF: %s", err)
		}
		anchorBTC.balances[i].wif = wif
		log.Debug("balance \n", i, spew.Sdump(anchorBTC.balances[i]))
	}

	//time.Sleep(1 * time.Second)
	return nil
}

func (anchorBTC *AnchorBTC) createBtcwalletNotificationHandlers() btcrpcclient.NotificationHandlers {

	ntfnHandlers := btcrpcclient.NotificationHandlers{
		OnAccountBalance: func(account string, balance btcutil.Amount, confirmed bool) {
			//go newBalance(account, balance, confirmed)
			log.Info("wclient: OnAccountBalance, account=", account, ", balance=",
				balance.ToUnit(btcutil.AmountBTC), ", confirmed=", confirmed)
		},

		OnWalletLockState: func(locked bool) {
			log.Info("wclient: OnWalletLockState, locked=", locked)
			anchorBTC.walletLocked = locked
		},

		OnUnknownNotification: func(method string, params []json.RawMessage) {
			log.Info("wclient: OnUnknownNotification: method=", method, "\nparams[0]=",
				string(params[0]), "\nparam[1]=", string(params[1]))
		},
	}

	return ntfnHandlers
}

func (anchorBTC *AnchorBTC) createBtcdNotificationHandlers() btcrpcclient.NotificationHandlers {

	ntfnHandlers := btcrpcclient.NotificationHandlers{

		OnBlockConnected: func(hash *chainhash.Hash, height int32, t time.Time) {
			log.Info("dclient: OnBlockConnected: hash=", hash, ", height=", height, ", time=", t)
			//go newBlock(hash, height)	// no need
		},

		OnRecvTx: func(transaction *btcutil.Tx, details *btcjson.BlockDetails) {
			log.Info("dclient: OnRecvTx: details=%#v\n", details)
			log.Info("dclient: OnRecvTx: tx=%#v,  tx.Hash=%#v, tx.index=%d\n",
				transaction, transaction.Hash().String(), transaction.Index())
		},

		OnRedeemingTx: func(transaction *btcutil.Tx, details *btcjson.BlockDetails) {
			log.Info("dclient: OnRedeemingTx: details=", details)
			log.Info("dclient: OnRedeemingTx: tx.Hash=", transaction.Hash().String(), "tx.index=", transaction.Index())

			if details != nil {
				// do not block OnRedeemingTx callback
				log.Info("Anchor: saveAnchorEntryInfo.")
				go anchorBTC.saveAnchorEntryInfo(transaction, details)
			}
		},

		OnRelevantTxAccepted: func(transaction []byte) {
			log.Info("dclient: OnRelevantTxAccepted: tx=", string(transaction))
		},
	}

	return ntfnHandlers
}

func (anchorBTC *AnchorBTC) saveAnchorEntryInfo(transaction *btcutil.Tx, details *btcjson.BlockDetails) {
	log.Info("in saveAnchorEntryInfo, anchor record count ", len(anchorBTC.toAnchorInfo))
	var saved = false
	for _, anchorInfo := range anchorBTC.toAnchorInfo {
		if strings.Compare(anchorInfo.Bitcoin.TXID, transaction.Hash().String()) == 0 {
			log.Info("Got the transaction return from bitcoin network")
			anchorInfo.Bitcoin.Address = anchorBTC.defaultAddress.String()
			anchorInfo.Bitcoin.BlockHeight = details.Height
			anchorInfo.Bitcoin.BlockHash = details.Hash
			anchorInfo.Bitcoin.Offset = int32(details.Index)
			log.Info("anchor.record saved: " + spew.Sdump(anchorInfo))

			err := anchorBTC.service.submitEntryToAnchorChain(anchorInfo)
			if err != nil {
				log.Error("Error in writing anchor into anchor chain: ", err)

			} else {
				saved = true
			}

			delete(anchorBTC.toAnchorInfo, anchorInfo.KeyMR)
			break
		}
	}
	// This happends when there's a double spending (for dir block 122 and its btc tx)
	// (see https://www.blocktrail.com/BTC/tx/ac82f4173259494b22f4987f1e18608f38f1ff756fb4a3c637dfb5565aa5e6cf)
	// or tx mutation / malleated
	// In this case, it will end up being re-anchored.
	if !saved {
		log.Info("Not saved to db: btc.tx=%s\n blockDetails=%s\n", spew.Sdump(transaction), spew.Sdump(details))
	}
}

func (anchorBTC *AnchorBTC) createRawTransaction(b balance, hash []byte, blockHeight uint32) (*wire.MsgTx, error) {
	msgtx := wire.NewMsgTx(wire.TxVersion)

	if err := anchorBTC.addTxOuts(msgtx, b, hash, blockHeight); err != nil {
		return nil, fmt.Errorf("cannot addTxOuts: %s", err)
	}

	if err := anchorBTC.addTxIn(msgtx, b); err != nil {
		return nil, fmt.Errorf("cannot addTxIn: %s", err)
	}

	if err := validateMsgTx(msgtx, []btcjson.ListUnspentResult{b.unspentResult}); err != nil {
		return nil, fmt.Errorf("cannot validateMsgTx: %s", err)
	}

	return msgtx, nil
}

func validateMsgTx(msgtx *wire.MsgTx, inputs []btcjson.ListUnspentResult) error {
	flags := txscript.ScriptBip16 | txscript.ScriptStrictMultiSig //ScriptCanonicalSignatures
	bip16 := time.Now().After(txscript.Bip16Activation)
	if bip16 {
		flags |= txscript.ScriptBip16
	}

	for i := range msgtx.TxIn {
		scriptPubKey, err := hex.DecodeString(inputs[i].ScriptPubKey)
		if err != nil {
			return fmt.Errorf("cannot decode scriptPubKey: %s", err)
		}
		engine, err := txscript.NewEngine(scriptPubKey, msgtx, i, flags, nil)
		if err != nil {
			log.Errorf("cannot create script engine: %s\n", err)
			return fmt.Errorf("cannot create script engine: %s", err)
		}
		if err = engine.Execute(); err != nil {
			log.Errorf("cannot execute script engine: %s\n  === UnspentResult: %s", err, spew.Sdump(inputs[i]))
			return fmt.Errorf("cannot execute script engine: %s", err)
		}
	}
	return nil
}

func (anchorBTC *AnchorBTC) addTxIn(msgtx *wire.MsgTx, b balance) error {
	output := b.unspentResult
	//anchorLog.Infof("unspentResult: %s\n", spew.Sdump(output))
	prevTxHash, err := chainhash.NewHashFromStr(output.TxID)
	if err != nil {
		return fmt.Errorf("cannot get sha hash from str: %s", err)
	}
	if prevTxHash == nil {
		log.Error("prevTxHash == nil")
	}

	outPoint := wire.NewOutPoint(prevTxHash, output.Vout)
	msgtx.AddTxIn(wire.NewTxIn(outPoint, nil))
	if outPoint == nil {
		log.Error("outPoint == nil")
	}

	// OnRedeemingTx
	err = anchorBTC.dclient.NotifySpent([]*wire.OutPoint{outPoint})
	if err != nil {
		log.Error("NotifySpent err: ", err)
	}

	subscript, err := hex.DecodeString(output.ScriptPubKey)
	if err != nil {
		return fmt.Errorf("cannot decode scriptPubKey: %s", err)
	}
	if subscript == nil {
		log.Error("subscript == nil")
	}

	sigScript, err := txscript.SignatureScript(msgtx, 0, subscript, txscript.SigHashAll, b.wif.PrivKey, true)
	if err != nil {
		return fmt.Errorf("cannot create scriptSig: %s", err)
	}
	if sigScript == nil {
		log.Error("sigScript == nil")
	}

	msgtx.TxIn[0].SignatureScript = sigScript
	return nil
}

func (anchorBTC *AnchorBTC) addTxOuts(msgtx *wire.MsgTx, b balance, hash []byte, blockHeight uint32) error {
	anchorHash, err := prependBlockHeight(blockHeight, hash)
	if err != nil {
		log.Errorf("ScriptBuilder error: %v\n", err)
	}

	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_RETURN)
	builder.AddData(anchorHash)

	opReturn, err := builder.Script()
	msgtx.AddTxOut(wire.NewTxOut(0, opReturn))
	if err != nil {
		log.Errorf("ScriptBuilder error: %v\n", err)
	}

	amount, _ := btcutil.NewAmount(b.unspentResult.Amount)
	change := amount - anchorBTC.fee

	// Check if there are leftover unspent outputs, and return coins back to
	// a new address we own.
	if change > 0 {

		// Spend change.
		pkScript, err := getPayScript(b.address)
		if err != nil {
			log.Error("PayToAddrScript error, address is ", spew.Sdump(b.address))
			return fmt.Errorf("cannot create txout script: %s", err)
		}
		msgtx.AddTxOut(wire.NewTxOut(int64(change), pkScript))
	}
	return nil
}

func getPayScript(address btcutil.Address) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
		AddData(address.ScriptAddress()).AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).
		Script()
}

func (anchorBTC *AnchorBTC) sendRawTransaction(msgtx *wire.MsgTx) (*chainhash.Hash, error) {
	//anchorLog.Debug("sendRawTransaction: msgTx=", spew.Sdump(msgtx))
	buf := bytes.Buffer{}
	buf.Grow(msgtx.SerializeSize())
	if err := msgtx.BtcEncode(&buf, wire.ProtocolVersion); err != nil {
		return nil, err
	}

	log.Debug("send transaction: ", spew.Sdump(msgtx))
	// use rpc client for btcd here for better callback info
	// this should not require wallet to be unlocked
	shaHash, err := anchorBTC.dclient.SendRawTransaction(msgtx, false)
	if err != nil {
		return nil, fmt.Errorf("failed in rpcclient.SendRawTransaction: %s", err)
	}
	log.Info("btc txHash returned: ", shaHash) // new tx hash
	return shaHash, nil
}

func (anchorBTC *AnchorBTC) doTransaction(anchor *anchor.AnchorRecord, hash *common.Hash) {
	if len(anchorBTC.balances) == 0 {
		log.Warning("len(balances) == 0, start rescan UTXO *** ")
		anchorBTC.updateUTXO()
	}

	if len(anchorBTC.balances) == 0 {
		log.Warning("No balance in your wallet. No anchoring for now")
		return
	}

	b := anchorBTC.balances[0]
	anchorBTC.balances = anchorBTC.balances[1:]
	anchorBTC.defaultAddress = b.address

	log.Info("new balances.len=", len(anchorBTC.balances))

	msgtx, err := anchorBTC.createRawTransaction(b, hash.Bytes(), anchor.DBHeight)
	if err != nil {
		log.Error("cannot create Raw Transaction: ", err)
		log.Error("Abort do this transaction: ", spew.Sdump(anchor))
		return
	}

	shaHash, err := anchorBTC.sendRawTransaction(msgtx)
	if err != nil {
		log.Error("cannot send Raw Transaction: %s", err)
		log.Error("Abort do this transaction: ", spew.Sdump(msgtx), " anchor ", spew.Sdump(anchor))
		return
	}
	anchor.Bitcoin.TXID = shaHash.String()

	log.Info("New anchor transaction for :", anchor.Bitcoin.TXID)
	anchorBTC.toAnchorInfo[anchor.KeyMR] = anchor
}
