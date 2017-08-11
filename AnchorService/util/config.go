package util

import (
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gcfg.v1"
	"os"
	"os/user"
)

type AnchorServiceCfg struct {
	App struct {
		HomeDir    string
		FactomAddr string
		AnchorTo   int
	}

	Anchor struct {
		ServerECKey         string
		AnchorChainID       string
		SigKey              string
		ConfirmationsNeeded int
	}

	Btc struct {
		SendToBTCinSeconds int
		WalletPassphrase   string
		CertHomePath       string
		RpcClientHost      string
		RpcClientEndpoint  string
		RpcClientUser      string
		RpcClientPass      string
		BtcTransFee        float64
		CertHomePathBtcd   string
		RpcBtcdHost        string
		RpcUser            string
		RpcPass            string
	}

	Eth struct {
		AccountPassphrase string
		AccountAddress    string
		EthHttpHost       string
		Gas               string
		GasPrice          string
	}

	Log struct {
		LogPath  string
		LogLevel string
	}
}

// defaultConfig
const defaultConfig = `
; ------------------------------------------------------------------------------
; App settings
; ------------------------------------------------------------------------------
[app]
HomeDir								= ""
AnchorTo							= 0

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
AccountAddress                         = "0x1786726f7636ba45b4Bfe9Cb546D06C313150E4D"
AccountPassphrase                      = "Initial0"
EthHttpHost                            = "localhost:8545"
Gas                                    = "0x76c0"
GasPrice                               = "0x9184e72a000"

; ------------------------------------------------------------------------------
; logLevel - allowed values are: debug, info, warn, error, fatal, panic
; ------------------------------------------------------------------------------
[log]
logLevel 							= info
LogPath								= "Log"
`

var config *AnchorServiceCfg
var logger = log.WithFields(log.Fields{"module": "common"})

func ReadConfig() *AnchorServiceCfg {
	if config == nil {
		config = readFromLocalOrDefault()
		logger.Info("Got Config \n", spew.Sdump(config))
	}

	return config
}

func readFromLocalOrDefault() *AnchorServiceCfg {
	fileName := GetHomeDir() + "/.anchorservice/anchorservice.conf"
	cfg := new(AnchorServiceCfg)

	err := gcfg.ReadFileInto(cfg, fileName)
	if err != nil {
		logger.Info("Error on load file config, fileName = %s, err = %s", fileName, err)
		logger.Info("Use the default config ")
		err = gcfg.ReadStringInto(cfg, defaultConfig)
		if err != nil {
			panic(err)
		}
	}

	return cfg
}

func GetHomeDir() string {
	// Get the OS specific home directory via the Go standard lib.
	var homeDir string
	usr, err := user.Current()
	if err == nil {
		homeDir = usr.HomeDir
	}

	// Fall back to standard HOME environment variable that works
	// for most POSIX OSes if the directory from the Go standard
	// lib failed.
	if err != nil || homeDir == "" {
		homeDir = os.Getenv("HOME")
	}
	return homeDir
}
