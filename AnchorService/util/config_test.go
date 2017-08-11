package util_test

import (
	"github.com/compasses/GOProjects/AnchorService/util"
	"gopkg.in/gcfg.v1"
	"testing"
)

func TestLoadDefaultConfigFull(t *testing.T) {
	var defaultConfig string = `
	; ------------------------------------------------------------------------------
; App settings
; AnchorTo: 0 -- BTC, 1 -- ETH
; ------------------------------------------------------------------------------
[app]
HomeDir                               = ""
FactomAddr                            = "localhost:8088"
AnchorTo                              = 1

; ------------------------------------------------------------------------------
; anchor settings
; ------------------------------------------------------------------------------
[anchor]
ServerECKey                         = Es38XqZaMQmjtuLKK9c238QsU2vFsFR4StHTUM6fVZFsSejHWcui
AnchorChainID                       = df3ade9eec4b08d5379cc64270c30ea7315d8a8a1a69efe2b98a60ecdd69e604
SigKey                              = e89f1216745b61d056b5297be7bdecd7b82966feb2ae482f686e127bd1b2ff80
ConfirmationsNeeded                 = 1

; ------------------------------------------------------------------------------
; anchor to bitcoin network
; ------------------------------------------------------------------------------
[btc]
WalletPassphrase                      = "Initial0"
CertHomePath                          = "btcwallet"
RpcClientHost                         = "localhost:18332"
RpcClientEndpoint                     = "ws"
RpcClientUser                         = "jbi"
RpcClientPass                         = "jbi123456"
BtcTransFee                           = 0.001
CertHomePathBtcd                      = "btcd"
RpcBtcdHost                           = "localhost:18334"

; ------------------------------------------------------------------------------
; anchor to ethereum network
; ------------------------------------------------------------------------------
[eth]
AccountAddress                         = "0x100c8b406978a413c4305b3AA6074F734feE6C9c"
AccountPassphrase                      = "Initial0"
EthHttpHost                            = "localhost:8545"
GasPrice                               = "0x1"

; ------------------------------------------------------------------------------
; logLevel - allowed values are: debug, info, warn, error, fatal, panic
; ------------------------------------------------------------------------------
[log]
logLevel                              = debug
	`

	cfg := new(util.AnchorServiceCfg)
	gcfg.ReadStringInto(cfg, defaultConfig)
	t.Log("cfg file %q", cfg)
}
