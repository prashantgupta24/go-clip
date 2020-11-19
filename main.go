package main

import (
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

type customButton struct {
	widget.Button
	tapped bool
	id     int
}

func (b *customButton) Tapped(e *fyne.PointEvent) {
	// fmt.Println("clicked ...")
	if valToWrite, exists := btnTextMap[b.id]; exists {
		clip.WriteAll(valToWrite)
	}
	b.Disable()
	defer func() {
		time.Sleep(time.Millisecond * 30)
		b.Enable()
	}()
}

func newButton(id int) *customButton {
	b := &customButton{
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
	})
	w.SetFixedSize(true)

	var btnArray []*customButton

	box := widget.NewVBox()
	for i := 0; i < 10; i++ {
		button := newButton(i)
		box.Append(button)
		btnArray = append(btnArray, button)
	}
	w.SetContent(box)
	monitorClipboard(btnArray)
	w.ShowAndRun()

}

func monitorClipboard(btnArray []*customButton) {

	btnMap := make(map[string]int)

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
					// log.Printf("change received: '%s'", change)
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
				} else {
					log.Printf("channel has been closed. exiting..")
				}
			}
		}
	}()
}
