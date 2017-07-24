package main

import (
	"github.com/compasses/GOProjects/AnchorService/util"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"syscall"
)

var log = util.MainLogger

func main() {
	app := cli.NewApp()
	app.Name = "AnchorService"
	app.Action = func(c *cli.Context) error {
		return nil
	}

	log.Info("start app", app)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sig := <-sc
		log.Info("Got signal", 0, "signal", sig)
	}()
	app.Run(os.Args)
}
