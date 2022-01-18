package main

import (
	"SimpleQR/smartqr"
	"SimpleQR/startup"
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

	copyContent := systray.AddMenuItem("copy decoded content", "copy decoded content")
	encodeQR := systray.AddMenuItem("encode QR", "encode QR")
	systray.AddSeparator()
	_, autoStartUpMenuTitle, err := startup.CheckAutoStart()
	if err != nil {
		fmt.Println(err)
	}
	autoStartUp := systray.AddMenuItem(autoStartUpMenuTitle, autoStartUpMenuTitle)
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quits this app")
	systray.SetTooltip("Smart QR")

	go func() {
		for {
			select {
			case <-copyContent.ClickedCh:
				go smartqr.CopyContent()
			case <-encodeQR.ClickedCh:
				go smartqr.EncodeQR()
			case <-autoStartUp.ClickedCh:
				go func() {
					startUpState, _, err := startup.CheckAutoStart()
					if err != nil {
						fmt.Println(err)
					}
					if startUpState {
						err := startup.RemoveLink()
						if err != nil {
							fmt.Println(err)
						}
					} else {
						err := startup.MakeLink()
						if err != nil {
							fmt.Println(err)
						}
					}
					_, autoStartUpMenuTitle, err := startup.CheckAutoStart()
					if err == nil {
						autoStartUp.SetTitle(autoStartUpMenuTitle)
						autoStartUp.SetTooltip(autoStartUpMenuTitle)
					} else {
						fmt.Println(err)
					}
				}()
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
