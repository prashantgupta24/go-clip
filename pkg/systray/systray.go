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
		menuItemToVal:  make(map[*systray.MenuItem]string),
		valExistsMap:   make(map[string]bool),
		truncateLength: 25,
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
	configureMenu := systray.AddMenuItem("Configuration", "")
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

func clearSlots(menuItemArray []menuItem) {
	clipboardInstance.mutex.Lock()
	defer clipboardInstance.mutex.Unlock()

	for _, menuItem := range menuItemArray {
		menuItem.instance.SetTitle("")
		menuItem.instance.SetTooltip("")
		delete(clipboardInstance.valExistsMap, clipboardInstance.menuItemToVal[menuItem.instance])
		delete(clipboardInstance.menuItemToVal, menuItem.instance)
		clipboardInstance.nextMenuItemIndex = 0
	}
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
			menuItem := clipboardInstance.menuItemArray[i].instance
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
	clipboardInstance.mutex.Lock()
	defer clipboardInstance.mutex.Unlock()

	for i := 0; i < numSlots; i++ {
		systray.AddSeparator()
		menuItemInstance := systray.AddMenuItem("", "")
		menuItem := menuItem{
			instance: menuItemInstance,
		}

		//sub menu1
		subMenuObfuscate := menuItemInstance.AddSubMenuItem("Obfuscate Password", "")
		subMenuObfuscate.Hide()
		subMenuObfuscate.Disable()
		menuItem.subMenuItems = append(menuItem.subMenuItems, subMenuObfuscate)

		//sub menu2
		subMenuPinToggle := menuItemInstance.AddSubMenuItem("Pin item", "")
		subMenuPinToggle.Hide()
		subMenuPinToggle.Disable()
		menuItem.subMenuItems = append(menuItem.subMenuItems, subMenuPinToggle)

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
					newTitle.WriteString(val[:4])

					for i := 4; i < min(len(val), clipboardInstance.truncateLength); i++ {
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
					// for _, menuItem := range clipboardInstance.menuItemArray {
					for {
						menuItem := clipboardInstance.menuItemArray[clipboardInstance.nextMenuItemIndex]
						fmt.Println("Index : ", clipboardInstance.nextMenuItemIndex)

						if !menuItem.instance.Disabled() && !menuItem.instance.Checked() {
							fmt.Println("final : ", clipboardInstance.nextMenuItemIndex)
							//delete last entry, if exists
							delete(clipboardInstance.valExistsMap, clipboardInstance.menuItemToVal[menuItem.instance])
							delete(clipboardInstance.menuItemToVal, menuItem.instance)
							// if clipboardInstance.menuItemToVal[menuItem] == "" {
							clipboardInstance.valExistsMap[val] = true
							clipboardInstance.menuItemToVal[menuItem.instance] = val
							//truncate to fit on app
							valTrunc := val
							if len(val) > clipboardInstance.truncateLength {
								valTrunc = val[:clipboardInstance.truncateLength] + "... (" + strconv.Itoa(len(val)) + " chars)"
							}
							menuItem.instance.SetTitle(valTrunc)
							menuItem.instance.SetTooltip(val)
							for _, subMenuItem := range menuItem.subMenuItems {
								subMenuItem.Show()
								subMenuItem.Enable()
							}
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
