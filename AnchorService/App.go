package main

import (
	"AnchorService/anchor"
	"AnchorService/common"
	"AnchorService/util"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"syscall"
)

var log = util.MainLogger

func main() {
	app := cli.NewApp()
	app.Name = "AnchorService"

	BlockMsg := make(chan common.DirectoryBlockAnchorInfo, 100)

	app.Action = func(c *cli.Context) error {
		service := anchor.NewAnchorService(BlockMsg)
		factomSync := anchor.NewFactomSync(service)
		go service.Start()
		//go factomSync.StartSync()
		go factomSync.SyncUp()

		return nil
	}

	log.Info("AnchorService start...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sig := <-sc
		log.Info("Got signal", 0, "signal", sig)
		log.Info("Shut down gracefully ...")
		os.Exit(1)
	}()

	app.Run(os.Args)
	select {}
}
