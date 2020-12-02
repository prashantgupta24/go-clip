package systray

import (
	"strconv"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func truncateVal(clipboardInstance *clipboard, val string) string {
	valTrunc := val
	if len(val) > clipboardInstance.truncateLength {
		valTrunc = val[:clipboardInstance.truncateLength] + "... (" + strconv.Itoa(len(val)) + " chars)"
	}
	return valTrunc
}

func obfuscateVal(clipboardInstance *clipboard, menuItem menuItem) {
	val := clipboardInstance.menuItemToVal[menuItem.instance]
	var newTitle strings.Builder
	newTitle.WriteString(val[:min(len(val), clipboardInstance.pwShowLength)])

	for i := clipboardInstance.pwShowLength; i < min(len(val), clipboardInstance.truncateLength); i++ {
		newTitle.WriteString("*")
	}
	menuItem.instance.SetTitle(newTitle.String())
	menuItem.subMenuItems[obfuscateMenu].Check()
}

func deleteMenuItem(clipboardInstance *clipboard, menuItem menuItem) {
	menuItem.instance.SetTitle("")
	menuItem.instance.SetTooltip("")

	delete(clipboardInstance.valExistsMap, clipboardInstance.menuItemToVal[menuItem.instance])
	delete(clipboardInstance.menuItemToVal, menuItem.instance)

	for menuName, subMenu := range menuItem.subMenuItems {
		subMenu.Hide()
		subMenu.Disable()
		if menuName != pinMenu {
			subMenu.Uncheck()
		}
	}
}

func acceptVal(clipboardInstance *clipboard, menuItem menuItem, val string) {
	//truncate to fit on screen
	valTrunc := truncateVal(clipboardInstance, val)

	menuItem.instance.SetTitle(valTrunc)
	menuItem.instance.SetTooltip(val)

	clipboardInstance.valExistsMap[val] = true
	clipboardInstance.menuItemToVal[menuItem.instance] = val

	for _, subMenuItem := range menuItem.subMenuItems {
		subMenuItem.Show()
		subMenuItem.Enable()
	}
}

func exchangeMenuItems(clipboardInstance *clipboard, existingMenuItem, newMenuItem menuItem) {
	existingMenuItemVal := clipboardInstance.menuItemToVal[existingMenuItem.instance]
	newMenuItemVal := clipboardInstance.menuItemToVal[newMenuItem.instance]

	existingMenuObfuscateChecked := existingMenuItem.subMenuItems[obfuscateMenu].Checked()
	newMenuObfuscateChecked := newMenuItem.subMenuItems[obfuscateMenu].Checked()

	deleteMenuItem(clipboardInstance, existingMenuItem)
	deleteMenuItem(clipboardInstance, newMenuItem)

	acceptVal(clipboardInstance, existingMenuItem, newMenuItemVal)
	acceptVal(clipboardInstance, newMenuItem, existingMenuItemVal)

	if existingMenuObfuscateChecked {
		obfuscateVal(clipboardInstance, newMenuItem)
	}
	if newMenuObfuscateChecked {
		obfuscateVal(clipboardInstance, existingMenuItem)
	}
}

func substituteMenuItem(clipboardInstance *clipboard, menuItem menuItem) {
	for i := 0; i < clipboardInstance.activeSlots; i++ {
		existingMenuItem := clipboardInstance.menuItemArray[i]
		if !existingMenuItem.instance.Disabled() && !existingMenuItem.instance.Checked() {
			if existingMenuItem.instance != menuItem.instance {
				// fmt.Println("same")
				// return // same item
				exchangeMenuItems(clipboardInstance, existingMenuItem, menuItem)
			}
			//found the right menu item
			// fmt.Println("emenu item : ", existingMenuItem.instance)
			// fmt.Println("menu item : ", menuItem.instance)

			// temp := existingMenuItem.instance
			// existingMenuItem.instance = menuItem.instance
			// existingMenuItem.instance.Enable()
			// menuItem.instance.Enable()
			// fmt.Println("menu item : ", existingMenuItem.instance)
			// menuItem.instance = temp
			// existingMenuItem.instance.SetTitle(clipboardInstance.menuItemToVal[menuItem.instance])
			// existingMenuItem.instance.SetTooltip(clipboardInstance.menuItemToVal[menuItem.instance])

			existingMenuItem.subMenuItems[pinMenu].Check()
			existingMenuItem.subMenuItems[pinMenu].SetTitle("Unpin item")
			// existingMenuItem.subMenuItems[0].SetTitle("Unpin item")
			// existingMenuItem.subMenuItems[0].Check()
			existingMenuItem.instance.Check()

			break
			// subMenuPinToggle.SetTitle("Unpin item")
			// subMenuPinToggle.Check()
			// menuItemInstance.Check()
		}
	}
}
