package main

import (
	"SimpleQR/smartqr"
	"fmt"
	"github.com/getlantern/systray"
	"io/ioutil"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	go smartqr.Run()

	systray.SetIcon(getIcon("icon/smartqr.ico"))

	encodeQR := systray.AddMenuItem("encode QR", "Local time")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quits this app")

	systray.SetTooltip("Smart QR")

	go func() {
		for {
			select {
			case <-encodeQR.ClickedCh:
				go smartqr.EncodeQR()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// Cleaning stuff here.
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
