package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/go-clip/clip"
)

var (
	btnTextMap map[int]string
)

func init() {
	btnTextMap = make(map[int]string)
}

type tappableButtonstruct struct {
	widget.Button
	tapped bool
	id     int
}

// func (b *tappableButtonstruct) CreateRenderer() fyne.WidgetRenderer {
// 	return widget.NewButton().CreateRenderer()
// }

func (b *tappableButtonstruct) Tapped(*fyne.PointEvent) {
	fmt.Println("here")
	clip.WriteAll(btnTextMap[b.id])
	b.tapped = true
	defer func() { // TODO move to a real animation
		time.Sleep(time.Millisecond * 500)
		b.tapped = false
		b.Refresh()
	}()
	b.Refresh()

	// if b.OnTapped != nil && !b.Disabled() {
	// 	b.OnTapped()
	// }
}

// func (t tappableButtonstruct) CreateRenderer() *widget.Button {

// }

func newButton(id int) *tappableButtonstruct {
	b := &tappableButtonstruct{
		id: id,
	}
	b.ExtendBaseWidget(b)
	return b
}

func main() {

	a := app.New()
	w := a.NewWindow("Clipboard")
	w.Resize(fyne.Size{
		Width: 400,
		// Height: 100,
	})
	w.SetFixedSize(true)

	// hello := widget.NewLabel("Clipboard")

	box := widget.NewVBox(
	// hello,
	)

	var btnArray []*tappableButtonstruct
	btnMap := make(map[string]int)
	// b1n := newButton()
	// b1n.SetText("hello")

	// b1 := widget.NewButton("Hello", func() {
	// 	fmt.Println("clicked...")
	// 	// clip.WriteAll()
	// })

	for i := 0; i < 10; i++ {
		// b1.Text = "asdf"
		b1n := newButton(i)
		box.Append(b1n)
		btnArray = append(btnArray, b1n)
	}
	// box.Append(b1n)
	w.SetContent(box)

	changes := make(chan string, 10)
	stopCh := make(chan struct{})

	go clip.Monitor(time.Second, stopCh, changes)

	// Watch for changes
	go func() {
		for {
			select {
			case <-stopCh:
				break
			default:
				change, ok := <-changes
				if ok {
					log.Printf("change received: '%s'", change)
					// b1.Text = change
					// b1.Refresh()
					val := strings.TrimSpace(change)
					if _, exists := btnMap[val]; !exists {
						for index, elem := range btnArray {
							if elem.Text == "" {
								btnMap[val] = index
								btnTextMap[elem.id] = val
								if len(val) > 20 {
									val = val[:20] + "... (" + strconv.Itoa(len(val)) + " chars)"
								}
								elem.Text = val
								elem.Refresh()
								break
							}
						}
					}

					// btnArray[0].Text = change
					// btnArray[0].Refresh()

				} else {
					log.Printf("channel has been closed. exiting..")
				}

			}
		}
	}()

	w.ShowAndRun()

}

// package main

// import (
// 	"image/color"

// 	"fyne.io/fyne"
// 	"fyne.io/fyne/app"
// 	"fyne.io/fyne/canvas"
// 	"fyne.io/fyne/layout"
// )

// func main() {
// 	myApp := app.New()
// 	myWindow := myApp.NewWindow("Box Layout")

// 	text1 := canvas.NewText("Hello", color.White)
// 	text2 := canvas.NewText("There", color.White)
// 	text3 := canvas.NewText("(right)", color.White)
// 	container := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
// 		text1, text2, layout.NewSpacer(), text3)

// 	text4 := canvas.NewText("centered", color.White)
// 	centered := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
// 		layout.NewSpacer(), text4, layout.NewSpacer())
// 	myWindow.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), container, centered))
// 	myWindow.ShowAndRun()
// }
