package main

import (
	"fmt"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/getlantern/systray"
	"github.com/go-clip/icon"
)

func main() {
	systray.Register(onReady, func() {})
	// app.Start()

	a := app.New()
	w := a.NewWindow("Hello")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(widget.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	w.ShowAndRun()
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
