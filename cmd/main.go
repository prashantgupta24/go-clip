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

type clipboard struct {
	menuItemArray     []*systray.MenuItem
	nextMenuItemIndex int
	menuItemToVal     map[*systray.MenuItem]string
	valExistsMap      map[string]bool
	activeSlots       int
}

var clipboardInstance *clipboard

func init() {
	clipboardInstance = &clipboard{
		menuItemToVal: make(map[*systray.MenuItem]string),
		activeSlots:   20,
		valExistsMap:  make(map[string]bool),
	}
}

func main() {
	systray.Run(onReady, func() {})
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTooltip("Clipboard")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	// We can manipulate the systray in other goroutines
	go func() {
		configureMenu := systray.AddMenuItem("Configuration", "Configuration")
		slotsMenu := configureMenu.AddSubMenuItem("slotsMenu", "SubMenu Test (middle)")
		slots5 := slotsMenu.AddSubMenuItem("5", "5")
		slots10 := slotsMenu.AddSubMenuItem("10", "10")
		slots20 := slotsMenu.AddSubMenuItem("20", "20")
		clearMenu := configureMenu.AddSubMenuItem("Clear", "Clear")

		addSlots(clipboardInstance.activeSlots, clipboardInstance)
		clipboardInstance.nextMenuItemIndex = 0
		changeActiveSlots(10, clipboardInstance)
		monitorClipboard(clipboardInstance)

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

		for {
			select {
			case <-slots5.ClickedCh:
				fmt.Println("changed to 5")
				changeActiveSlots(5, clipboardInstance)
			case <-slots10.ClickedCh:
				fmt.Println("changed to 10")
				changeActiveSlots(10, clipboardInstance)
			case <-slots20.ClickedCh:
				fmt.Println("changed to 20")
				changeActiveSlots(20, clipboardInstance)
			case <-clearMenu.ClickedCh:
				fmt.Println("clear")
				clearSlots(clipboardInstance.menuItemArray)
			}
		}
	}()
}

func clearSlots(menuItemArray []*systray.MenuItem) {
	for _, menuItem := range menuItemArray {
		menuItem.SetTitle("")
		menuItem.SetTooltip("")
		delete(clipboardInstance.valExistsMap, clipboardInstance.menuItemToVal[menuItem])
		delete(clipboardInstance.menuItemToVal, menuItem)
	}
}

func changeActiveSlots(changeSlotNumTo int, clipboardInstance *clipboard) {

	existingSlots := clipboardInstance.activeSlots
	clipboardInstance.activeSlots = changeSlotNumTo
	if changeSlotNumTo == existingSlots {
		return
	}
	if changeSlotNumTo > existingSlots { //enable
		for i := existingSlots; i < changeSlotNumTo; i++ {
			menuItem := clipboardInstance.menuItemArray[i]
			menuItem.Enable()
			menuItem.Show()
		}
		clipboardInstance.nextMenuItemIndex = existingSlots
	} else { //disable
		for i := changeSlotNumTo; i < existingSlots; i++ {
			menuItem := clipboardInstance.menuItemArray[i]
			menuItem.Disable()
			menuItem.Hide()
			menuItem.SetTitle("")
			menuItem.SetTooltip("")
			delete(clipboardInstance.valExistsMap, clipboardInstance.menuItemToVal[menuItem])
			delete(clipboardInstance.menuItemToVal, menuItem)
		}
		clipboardInstance.nextMenuItemIndex = 0
	}

}

func addSlots(numSlots int, clipboardInstance *clipboard) {
	for i := 0; i < numSlots; i++ {
		systray.AddSeparator()
		menuItem := systray.AddMenuItem("", "")
		clipboardInstance.menuItemArray = append(clipboardInstance.menuItemArray, menuItem)
		go func() {
			for {
				select {
				case <-menuItem.ClickedCh:
					if valToWrite, exists := clipboardInstance.menuItemToVal[menuItem]; exists {
						clip.WriteAll(valToWrite)
					}
				}
			}
		}()
	}
}

func monitorClipboard(clipboardInstance *clipboard) {

	changes := make(chan string, 10)
	stopCh := make(chan struct{})

	go clip.Monitor(time.Millisecond*500, stopCh, changes)

	// Watch for changes
	go func() {
		for {
			select {
			case <-stopCh:
				break
			default:
				change, ok := <-changes
				if ok {
					val := strings.TrimSpace(change)
					fmt.Println("val : ", val)

					if _, exists := clipboardInstance.valExistsMap[val]; val != "" && !exists {
						fmt.Println("clipboardInstance.nextMenuItemIndex : ", clipboardInstance.nextMenuItemIndex)
						menuItem := clipboardInstance.menuItemArray[clipboardInstance.nextMenuItemIndex]
						// for _, menuItem := range clipboardInstance.menuItemArray {
						for {
							if !menuItem.Disabled() {
								//delete last entry, if exists
								delete(clipboardInstance.valExistsMap, clipboardInstance.menuItemToVal[menuItem])
								delete(clipboardInstance.menuItemToVal, menuItem)
								// if clipboardInstance.menuItemToVal[menuItem] == "" {
								clipboardInstance.valExistsMap[val] = true
								clipboardInstance.menuItemToVal[menuItem] = val
								//truncate to fit on app
								valTrunc := val
								if len(val) > 20 {
									valTrunc = val[:20] + "... (" + strconv.Itoa(len(val)) + " chars)"
								}
								menuItem.SetTitle(valTrunc)
								menuItem.SetTooltip(val)
								clipboardInstance.nextMenuItemIndex = (clipboardInstance.nextMenuItemIndex + 1) % clipboardInstance.activeSlots
								break
							} else {
								// menuItem = clipboardInstance.menuItemArray[(clipboardInstance.nextMenuItemIndex+1)%(clipboardInstance.activeSlots)]
								clipboardInstance.nextMenuItemIndex = (clipboardInstance.nextMenuItemIndex + 1) % clipboardInstance.activeSlots
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
