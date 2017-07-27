package util

import (
	"fmt"
	"gopkg.in/gcfg.v1"
	"os"
	"os/user"
)

type AnchorServiceCfg struct {
	App struct {
		HomeDir  string
		AnchorTo int
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
		WalletPassphrase  string
		CertHomePath      string
		RpcClientHost     string
		RpcClientEndpoint string
		RpcClientUser     string
		RpcClientPass     string
		EthTransFee       float64
		CertHomePathBtcd  string
		RpcEthHost        string
		RpcUser           string
		RpcPass           string
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
; logLevel - allowed values are: debug, info, warn, error, fatal, panic
; ------------------------------------------------------------------------------
[log]
logLevel 							= info
LogPath								= "Log"
`

var config *AnchorServiceCfg

func ReadConfig() *AnchorServiceCfg {
	if config == nil {
		config = readFromLocalOrDefault()
		fmt.Printf("Got Config %q\n", config)
	}

	return config
}

func readFromLocalOrDefault() *AnchorServiceCfg {
	fileName := getHomeDir() + "/.anchorservice/anchorservice.conf"
	cfg := new(AnchorServiceCfg)

	err := gcfg.ReadFileInto(cfg, fileName)
	if err != nil {
		fmt.Errorf("Error on load file config, fileName = %s, err = %s", fileName, err)
		fmt.Println("Use the default config ")
		err = gcfg.ReadStringInto(cfg, defaultConfig)
		if err != nil {
			panic(err)
		}
	}

	return cfg
}

func getHomeDir() string {
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
