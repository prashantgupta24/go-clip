package systray

import (
	"strconv"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func deleteMenuItem(clipboardInstance *clipboard, menuItem menuItem) {
	menuItem.instance.SetTitle("")
	menuItem.instance.SetTooltip("")
	delete(clipboardInstance.valExistsMap, clipboardInstance.menuItemToVal[menuItem.instance])
	delete(clipboardInstance.menuItemToVal, menuItem.instance)

}

func acceptVal(clipboardInstance *clipboard, menuItem menuItem, val string) {

	clipboardInstance.valExistsMap[val] = true
	clipboardInstance.menuItemToVal[menuItem.instance] = val

	//truncate to fit on screen
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
}
