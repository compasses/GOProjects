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
; ------------------------------------------------------------------------------
[app]
HomeDir								= ""

[anchor]
ServerECKey							= 397c49e182caa97737c6b394591c614156fbe7998d7bf5d76273961e9fa1edd406ed9e69bfdf85db8aa69820f348d096985bc0b11cc9fc9dcee3b8c68b41dfd5
AnchorChainID						= df3ade9eec4b08d5379cc64270c30ea7315d8a8a1a69efe2b98a60ecdd69e604
ConfirmationsNeeded					= 20

[btc]
WalletPassphrase 	  				= "lindasilva"
CertHomePath			  			= "btcwallet"
RpcClientHost			  			= "localhost:18332"
RpcClientEndpoint					= "ws"
RpcClientUser			  			= "testuser"
RpcClientPass 						= "notarychain"
BtcTransFee				  			= 0.0001
CertHomePathBtcd					= "btcd"
RpcBtcdHost 			  			= "localhost:18334"
RpcUser								= testuser
RpcPass								= notarychain

[eth]
WalletPassphrase 	  				= "lindasilva"
CertHomePath			  			= "btcwallet"
RpcClientHost			  			= "localhost:18332"
RpcClientEndpoint					= "ws"
RpcClientUser			  			= "testuser"
RpcClientPass 						= "notarychain"
EthTransFee				  			= 0.0001
CertHomePathBtcd					= "btcd"
RpcEthHost 			  			= "localhost:18334"
RpcUser								= testuser
RpcPass								= notarychain

; ------------------------------------------------------------------------------
; logLevel - allowed values are: debug, info, notice, warning, error, critical, alert, emergency and none
; ------------------------------------------------------------------------------
[log]
logLevel 							= info
LogPath								= "Log"
	`

	cfg := new(util.AnchorServiceCfg)
	gcfg.ReadStringInto(cfg, defaultConfig)
	t.Log("cfg file %q", cfg)
}
