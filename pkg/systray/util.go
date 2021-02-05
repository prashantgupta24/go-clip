package systray

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func obfuscateVal(clipboardInstance *clipboard, menuItem menuItem) {
	val := clipboardInstance.menuItemToVal[menuItem.instance]
	var newTitle strings.Builder
	newTitle.WriteString(val[:min(len(val), clipboardInstance.pwShowLength)])

	for i := clipboardInstance.pwShowLength; i < min(len(val), clipboardInstance.truncateLength); i++ {
		newTitle.WriteString("*")
	}
	menuItem.instance.SetTitle(newTitle.String())
	menuItem.instance.SetTooltip(newTitle.String())
	menuItem.subMenuItems[obfuscateMenu].Check()
}

func acceptVal(clipboardInstance *clipboard, menuItem menuItem, val string) {
	//truncate to fit on screen
	valTrunc := truncateVal(clipboardInstance, val)

	// menuItem.instance.SetTitle(valTrunc)
	// menuItem.instance.SetTooltip(val)
	menuItem.instance.SetTitle(removeNonASCII(valTrunc))
	menuItem.instance.SetTooltip(removeNonASCII(val))

	clipboardInstance.valExistsMap[val] = true
	clipboardInstance.menuItemToVal[menuItem.instance] = val

	for _, subMenuItem := range menuItem.subMenuItems {
		subMenuItem.Show()
		subMenuItem.Enable()
	}
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
				exchangeMenuItems(clipboardInstance, existingMenuItem, menuItem)
			}
			existingMenuItem.subMenuItems[pinMenu].Check()
			existingMenuItem.subMenuItems[pinMenu].SetTitle("Unpin item")
			existingMenuItem.instance.Check()
			break
		}
	}
}

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

func removeNonASCII(val string) string {
	removeNonASCII := regexp.MustCompile("[[:^ascii:]]")
	nonASCIIStr := removeNonASCII.ReplaceAllLiteralString(val, "")
	return strings.TrimSpace(nonASCIIStr)
}

func getTitle(menuItem menuItem) string {
	title := ""
	menuItemReflect := reflect.ValueOf(menuItem.instance).Elem()

	for i := 0; i < menuItemReflect.NumField(); i++ {
		fieldName := menuItemReflect.Type().Field(i).Name
		if fieldName == "title" {
			title = menuItemReflect.Field(i).String()
		}
	}
	return title
}

func getToolTip(menuItem menuItem) string {
	toolTip := ""
	menuItemReflect := reflect.ValueOf(menuItem.instance).Elem()

	for i := 0; i < menuItemReflect.NumField(); i++ {
		fieldName := menuItemReflect.Type().Field(i).Name
		if fieldName == "tooltip" {
			toolTip = menuItemReflect.Field(i).String()
		}
	}
	return toolTip
}
