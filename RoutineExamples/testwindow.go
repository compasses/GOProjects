package main

import (
	"flag"

	"github.com/andlabs/ui"
	"github.com/cryptix/go/msgbox"
)

var (
	mainW ui.Window

	started chan struct{}
	newText chan string
)

func initGui() {
	read := ui.NewButton("Read")
	tf := ui.NewTextField()

	read.OnClicked(func() {
		newText <- tf.Text()
	})

	stack := ui.NewVerticalStack(
		read,
		tf,
	)
	stack.SetStretchy(1)

	mainW = ui.NewWindow("TestGui", 200, 300, stack)
	mainW.OnClosing(func() bool {
		ui.Stop()
		return true
	})
	mainW.Show()
	close(started)
}

func doStuff() {
	<-started

	for t := range newText {
		if len(t) > 0 && t != "nice" {
			msgbox.New(mainW, "Wrong input", "please enter nice")
		}
	}
}

func main() {
	flag.Parse()

	started = make(chan struct{})
	newText = make(chan string)
	go ui.Do(initGui)
	go doStuff()

	if err := ui.Go(); err != nil {
		panic(err)
	}
}
