package systray

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/go-clip/clip"
	"github.com/go-clip/icon"
)

var clipboardInstance *clipboard

func initInstance() {
	clipboardInstance = &clipboard{
		menuItemToVal:  make(map[*systray.MenuItem]string),
		valExistsMap:   make(map[string]bool),
		truncateLength: 25,
		pwShowLength:   4,
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

	configureMenu := systray.AddMenuItem("Configuration", "")
	slotsMenu := configureMenu.AddSubMenuItem("slotsMenu", "")
	slots5 := slotsMenu.AddSubMenuItem("5", "")
	slots10 := slotsMenu.AddSubMenuItem("10", "")
	slots20 := slotsMenu.AddSubMenuItem("20", "")
	clearMenu := configureMenu.AddSubMenuItem("Clear", "Clear all entries (except pinned)")

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
			changeActiveSlots(5, clipboardInstance)
		case <-slots10.ClickedCh:
			changeActiveSlots(10, clipboardInstance)
		case <-slots20.ClickedCh:
			changeActiveSlots(20, clipboardInstance)
		case <-clearMenu.ClickedCh:
			clearSlots(clipboardInstance.menuItemArray)
		}
	}
}

func clearSlots(menuItemArray []menuItem) {
	clipboardInstance.mutex.Lock()
	defer clipboardInstance.mutex.Unlock()

	for _, menuItem := range menuItemArray {
		if !menuItem.instance.Checked() {
			deleteMenuItem(clipboardInstance, menuItem)
		}
	}
	clipboardInstance.nextMenuItemIndex = 0
}

func changeActiveSlots(changeSlotNumTo int, clipboardInstance *clipboard) {
	clipboardInstance.mutex.Lock()
	defer clipboardInstance.mutex.Unlock()

	existingSlots := clipboardInstance.activeSlots
	clipboardInstance.activeSlots = changeSlotNumTo
	if changeSlotNumTo == existingSlots {
		return
	}
	if changeSlotNumTo > existingSlots { //enable
		for i := existingSlots; i < changeSlotNumTo; i++ {
			menuItem := clipboardInstance.menuItemArray[i].instance
			menuItem.Enable()
			menuItem.Show()
		}
		for index, menuItem := range clipboardInstance.menuItemArray {
			if _, exists := clipboardInstance.menuItemToVal[menuItem.instance]; !exists && !menuItem.instance.Disabled() {
				clipboardInstance.nextMenuItemIndex = index
				break
			}
		}
	} else { //disable
		for i := changeSlotNumTo; i < existingSlots; i++ {
			menuItem := clipboardInstance.menuItemArray[i]
			menuItem.instance.Disable()
			menuItem.instance.Hide()
			deleteMenuItem(clipboardInstance, menuItem)
		}
		if clipboardInstance.nextMenuItemIndex >= changeSlotNumTo {
			clipboardInstance.nextMenuItemIndex = 0
		}
	}

}

func addSlots(numSlots int, clipboardInstance *clipboard) {
	clipboardInstance.mutex.Lock()
	defer clipboardInstance.mutex.Unlock()

	for i := 0; i < numSlots; i++ {
		systray.AddSeparator()
		menuItemInstance := systray.AddMenuItem("", "")
		menuItem := menuItem{
			instance: menuItemInstance,
		}

		//sub menu1
		subMenuPinToggle := menuItemInstance.AddSubMenuItem("Pin item", "")
		subMenuPinToggle.Hide()
		subMenuPinToggle.Disable()
		menuItem.subMenuItems = append(menuItem.subMenuItems, subMenuPinToggle)

		//sub menu2
		subMenuObfuscate := menuItemInstance.AddSubMenuItem("Obfuscate Password", "")
		subMenuObfuscate.Hide()
		subMenuObfuscate.Disable()
		menuItem.subMenuItems = append(menuItem.subMenuItems, subMenuObfuscate)

		clipboardInstance.menuItemArray = append(clipboardInstance.menuItemArray, menuItem)
		go func() {
			for {
				select {
				case <-menuItemInstance.ClickedCh:
					if valToWrite, exists := clipboardInstance.menuItemToVal[menuItemInstance]; exists {
						clip.WriteAll(valToWrite)
					}
				case <-subMenuObfuscate.ClickedCh:
					val := clipboardInstance.menuItemToVal[menuItemInstance]
					var newTitle strings.Builder
					newTitle.WriteString(val[:clipboardInstance.pwShowLength])

					for i := clipboardInstance.pwShowLength; i < min(len(val), clipboardInstance.truncateLength); i++ {
						newTitle.WriteString("*")
					}
					menuItemInstance.SetTitle(newTitle.String())

				case <-subMenuPinToggle.ClickedCh:
					if subMenuPinToggle.Checked() {
						subMenuPinToggle.SetTitle("Pin item")
						subMenuPinToggle.Uncheck()
						menuItemInstance.Uncheck()
					} else {
						subMenuPinToggle.SetTitle("Unpin item")
						subMenuPinToggle.Check()
						menuItemInstance.Check()
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
					for {
						menuItem := clipboardInstance.menuItemArray[clipboardInstance.nextMenuItemIndex]
						// fmt.Println("Index : ", clipboardInstance.nextMenuItemIndex)
						if !menuItem.instance.Disabled() && !menuItem.instance.Checked() {
							// fmt.Println("final : ", clipboardInstance.nextMenuItemIndex)
							deleteMenuItem(clipboardInstance, menuItem) //delete last entry, if exists
							acceptVal(clipboardInstance, menuItem, val) //add new entry
							clipboardInstance.nextMenuItemIndex = (clipboardInstance.nextMenuItemIndex + 1) % clipboardInstance.activeSlots
							break
						} else {
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
