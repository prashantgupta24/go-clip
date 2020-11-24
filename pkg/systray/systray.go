package systray

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

var clipboardInstance *clipboard

func initInstance() {
	clipboardInstance = &clipboard{
		menuItemToVal: make(map[*systray.MenuItem]string),
		valExistsMap:  make(map[string]bool),
	}
}

//Run starts the system tray app
func Run() {
	initInstance()
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
	configureMenu := systray.AddMenuItem("Configuration", "Configuration")
	slotsMenu := configureMenu.AddSubMenuItem("slotsMenu", "SubMenu Test (middle)")
	slots5 := slotsMenu.AddSubMenuItem("5", "5")
	slots10 := slotsMenu.AddSubMenuItem("10", "10")
	slots20 := slotsMenu.AddSubMenuItem("20", "20")
	clearMenu := configureMenu.AddSubMenuItem("Clear", "Clear")

	addSlots(20, clipboardInstance)
	changeActiveSlots(10, clipboardInstance)

	//monitor clipboard
	changes := make(chan string, 10)
	stopCh := make(chan struct{})
	go clip.Monitor(time.Millisecond*500, stopCh, changes)
	go monitorClipboard(clipboardInstance, stopCh, changes)

	for {
		select {
		case <-slots5.ClickedCh:
			// fmt.Println("changed to 5")
			changeActiveSlots(5, clipboardInstance)
		case <-slots10.ClickedCh:
			// fmt.Println("changed to 10")
			changeActiveSlots(10, clipboardInstance)
		case <-slots20.ClickedCh:
			// fmt.Println("changed to 20")
			changeActiveSlots(20, clipboardInstance)
		case <-clearMenu.ClickedCh:
			// fmt.Println("clear")
			clearSlots(clipboardInstance.menuItemArray)
		}
	}
}

func clearSlots(menuItemArray []*systray.MenuItem) {
	for _, menuItem := range menuItemArray {
		menuItem.SetTitle("")
		menuItem.SetTooltip("")
		delete(clipboardInstance.valExistsMap, clipboardInstance.menuItemToVal[menuItem])
		delete(clipboardInstance.menuItemToVal, menuItem)
		clipboardInstance.nextMenuItemIndex = 0
	}
}

func changeActiveSlots(changeSlotNumTo int, clipboardInstance *clipboard) {
	clipboardInstance.mutex.RLock()
	defer func() {
		// fmt.Println("release Rlock")
		clipboardInstance.mutex.RUnlock()
	}()

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
		for index, menuItem := range clipboardInstance.menuItemArray {
			if _, exists := clipboardInstance.menuItemToVal[menuItem]; !exists && !menuItem.Disabled() {
				clipboardInstance.nextMenuItemIndex = index
				break
			}
		}
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
		if clipboardInstance.nextMenuItemIndex >= changeSlotNumTo {
			clipboardInstance.nextMenuItemIndex = 0
		}
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
	clipboardInstance.activeSlots = clipboardInstance.activeSlots + numSlots
}

func monitorClipboard(clipboardInstance *clipboard, stopCh chan struct{}, changes chan string) {
	// Watch for changes
	for {
		select {
		case <-stopCh:
			break
		default:
			change, ok := <-changes
			if ok {
				clipboardInstance.mutex.Lock()

				val := strings.TrimSpace(change)
				// fmt.Println("val : ", val)

				if _, exists := clipboardInstance.valExistsMap[val]; val != "" && !exists {
					// fmt.Println("Index : ", clipboardInstance.nextMenuItemIndex)
					menuItem := clipboardInstance.menuItemArray[clipboardInstance.nextMenuItemIndex]
					// for _, menuItem := range clipboardInstance.menuItemArray {
					for {
						if !menuItem.Disabled() {
							// fmt.Println("final : ", clipboardInstance.nextMenuItemIndex)
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
				// fmt.Println("release lock")
				clipboardInstance.mutex.Unlock()
			} else {
				log.Printf("channel has been closed. exiting..")
			}
		}
	}
}
