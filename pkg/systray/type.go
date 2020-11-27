package systray

import (
	"sync"

	"github.com/getlantern/systray"
)

type clipboard struct {
	menuItemArray     []menuItem
	nextMenuItemIndex int
	menuItemToVal     map[*systray.MenuItem]string
	valExistsMap      map[string]bool
	activeSlots       int
	truncateLength    int
	mutex             sync.RWMutex
}

type menuItem struct {
	instance     *systray.MenuItem
	subMenuItems []*systray.MenuItem
}
