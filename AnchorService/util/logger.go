package util

import (
	log "github.com/sirupsen/logrus"
)

var ContextLogger *log.Entry
var MainLogger *log.Entry
var CommonLogger *log.Entry
var AnchorLogger *log.Entry

func init() {
	cfg := ReadConfig()
	log.SetLevel(GetLogLevel(cfg))
	ContextLogger = log.WithFields(log.Fields{"common": "anchorservice"})

	MainLogger = log.WithFields(log.Fields{"module": "main"})
	CommonLogger = log.WithFields(log.Fields{"module": "common"})
	AnchorLogger = log.WithFields(log.Fields{"module": "anchor"})
}

func GetLogLevel(cfg *AnchorServiceCfg) log.Level {
	switch cfg.Log.LogLevel {
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "debug":
		return log.DebugLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		panic("No known log level" + cfg.Log.LogLevel)
	}
}
