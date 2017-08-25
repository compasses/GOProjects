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
		Interval            int
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
; AnchorTo: 0 -- BTC, 1 -- ETH
; ------------------------------------------------------------------------------
[app]
HomeDir                               = ""
FactomAddr                            = "localhost:8088"
AnchorTo                              = 1

; ------------------------------------------------------------------------------
; anchor settings
; factom related
; Interval (minutes)
; ------------------------------------------------------------------------------
[anchor]
ServerECKey                         = Es38XqZaMQmjtuLKK9c238QsU2vFsFR4StHTUM6fVZFsSejHWcui
AnchorChainID                       = df3ade9eec4b08d5379cc64270c30ea7315d8a8a1a69efe2b98a60ecdd69e604
SigKey                              = e89f1216745b61d056b5297be7bdecd7b82966feb2ae482f686e127bd1b2ff80
Interval                            = 10
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
