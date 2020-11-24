package main

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/go-clip/icon"
	"github.com/go-clip/pkg/app"
)

func main() {
	systray.Register(onReady, func() {})
	app.Start()
}

func onReady() {
	fmt.Println("here??")
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Webview example")
	mShowGoogle := systray.AddMenuItem("Show Google", "")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		for {
			select {
			case <-mShowGoogle.ClickedCh:
				fmt.Println("clicked")
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()

}
