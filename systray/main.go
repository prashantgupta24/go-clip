package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/go-clip/clip"
	"github.com/go-clip/icon"
)

var (
	btnTextMap map[*systray.MenuItem]string
)

func init() {
	btnTextMap = make(map[*systray.MenuItem]string)
}

func main() {
	systray.Run(onReady, func() {})
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	// systray.SetTitle("Awesome App")
	systray.SetTooltip("Clipboard")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	// We can manipulate the systray in other goroutines
	go func() {
		// systray.SetTemplateIcon(icon.Data, icon.Data)
		// systray.SetTitle("Awesome App")
		// systray.SetTooltip("Pretty awesome棒棒嗒")

		// mChange := systray.AddMenuItem("Change Me", "Change Me")
		// mChecked := systray.AddMenuItemCheckbox("Unchecked", "Check Me", true)
		// mEnabled := systray.AddMenuItem("Enabled", "Enabled")
		// // Sets the icon of a menu item. Only available on Mac.
		// mEnabled.SetTemplateIcon(icon.Data, icon.Data)

		// systray.AddMenuItem("Ignored", "Ignored")

		subMenuTop := systray.AddMenuItem("SubMenuTop", "SubMenu Test (top)")
		subMenuMiddle := subMenuTop.AddSubMenuItem("SubMenuMiddle", "SubMenu Test (middle)")
		subMenuMiddle.AddSubMenuItemCheckbox("SubMenuBottom - Toggle Panic!", "SubMenu Test (bottom) - Hide/Show Panic!", false)
		subMenuMiddle.AddSubMenuItem("SubMenuBottom - Panic!", "SubMenu Test (bottom)")

		// mUrl := systray.AddMenuItem("Open UI", "my home")
		// mQuit := systray.AddMenuItem("退出", "Quit the whole app")

		// // Sets the icon of a menu item. Only available on Mac.
		// mQuit.SetIcon(icon.Data)

		// systray.AddSeparator()
		var menuItemArray []*systray.MenuItem
		for i := 0; i < 10; i++ {
			systray.AddSeparator()
			menuItem := systray.AddMenuItem("", "")
			menuItemArray = append(menuItemArray, menuItem)
			go func() {
				for {
					select {
					case <-menuItem.ClickedCh:
						// a := menuItem.String()
						// fmt.Println("a : ", menuItem)
						if valToWrite, exists := btnTextMap[menuItem]; exists {
							clip.WriteAll(valToWrite)
						}
						// menuItem.SetTitle("I've Changed")
					}
				}

			}()
		}

		monitorClipboard(menuItemArray)

		// mToggle := systray.AddMenuItem("Toggle", "Toggle the Quit button")
		// shown := true
		// toggle := func() {
		// 	if shown {
		// 		subMenuBottom.Check()
		// 		subMenuBottom2.Hide()
		// 		mQuitOrig.Hide()
		// 		mEnabled.Hide()
		// 		shown = false
		// 	} else {
		// 		subMenuBottom.Uncheck()
		// 		subMenuBottom2.Show()
		// 		mQuitOrig.Show()
		// 		mEnabled.Show()
		// 		shown = true
		// 	}
		// }

		// for {
		// 	select {
		// 	case <-mChange.ClickedCh:
		// 		mChange.SetTitle("I've Changed")
		// 	case <-mChecked.ClickedCh:
		// 		if mChecked.Checked() {
		// 			mChecked.Uncheck()
		// 			mChecked.SetTitle("Unchecked")
		// 		} else {
		// 			mChecked.Check()
		// 			mChecked.SetTitle("Checked")
		// 		}
		// 	case <-mEnabled.ClickedCh:
		// 		mEnabled.SetTitle("Disabled")
		// 		mEnabled.Disable()
		// 	case <-mUrl.ClickedCh:
		// 		open.Run("https://www.getlantern.org")
		// 	case <-subMenuBottom2.ClickedCh:
		// 		panic("panic button pressed")
		// 	case <-subMenuBottom.ClickedCh:
		// 		toggle()
		// 	case <-mToggle.ClickedCh:
		// 		toggle()
		// 	case <-mQuit.ClickedCh:
		// 		systray.Quit()
		// 		fmt.Println("Quit2 now...")
		// 		return
		// 	}
		// }
	}()
}

func monitorClipboard(menuItemArray []*systray.MenuItem) {

	btnMap := make(map[string]bool)

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
					// fmt.Println("val : ", val)

					if _, exists := btnMap[val]; !exists {

						for _, menuItem := range menuItemArray {
							if btnTextMap[menuItem] == "" {
								btnMap[val] = true
								valTrunc := val
								if len(val) > 20 {
									valTrunc = val[:20] + "... (" + strconv.Itoa(len(val)) + " chars)"
								}
								menuItem.SetTitle(valTrunc)
								btnTextMap[menuItem] = val
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
