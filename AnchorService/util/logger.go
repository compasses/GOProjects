package util

import (
	log "github.com/inconshreveable/log15"
	"os"
)

var MainLogger log.Logger
var AnchorLogger log.Logger
var FactomLooger log.Logger

func init() {
	cfg := ReadConfig()
	logFile := GetHomeDir() + "/.anchorservice/anchorservice.log"
	var handler log.Handler

	if GetLogLevel(cfg) == log.LvlDebug {
		handler = log.MultiHandler(
			log.StreamHandler(os.Stdout, CustomFormat()),
		)
		log.Root().SetHandler(log.CallerFileHandler(handler))
	} else {
		handler = log.MultiHandler(
			log.Must.FileHandler(logFile, CustomFormat()),
		)

		log.Root().SetHandler(log.LvlFilterHandler(GetLogLevel(cfg), handler))
	}

	MainLogger = log.New("module", "main")
	AnchorLogger = log.New("module", "anchor")
	FactomLooger = log.New("module", "factom")
}

func GetLogLevel(cfg *AnchorServiceCfg) log.Lvl {
	switch cfg.Log.LogLevel {
	case "info":
		return log.LvlInfo
	case "warn":
		return log.LvlWarn
	case "debug":
		return log.LvlDebug
	case "error":
		return log.LvlError
	case "fatal":
		return log.LvlCrit
	default:
		panic("No known log level" + cfg.Log.LogLevel)
	}
}
