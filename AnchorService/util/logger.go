package util

import (
	log "github.com/inconshreveable/log15"
	"os"
)

var ContextLogger log.Logger
var MainLogger log.Logger
var CommonLogger log.Logger
var AnchorLogger log.Logger

func init() {
	cfg := ReadConfig()
	logFile := GetHomeDir() + "/.anchorservice/anchorservice.log"
	var handler log.Handler

	if GetLogLevel(cfg) == log.LvlDebug {
		handler = log.MultiHandler(
			log.StreamHandler(os.Stdout, log.LogfmtFormat()),
		)

		log.Root().SetHandler(log.CallerFileHandler(handler))
	} else {
		handler = log.MultiHandler(
			log.Must.FileHandler(logFile, log.LogfmtFormat()),
		)
		log.Root().SetHandler(handler)
	}

	ContextLogger = log.New("common", "anchorservice")

	MainLogger = log.New("module", "main")
	CommonLogger = log.New("module", "common")
	AnchorLogger = log.New("module", "anchor")
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
