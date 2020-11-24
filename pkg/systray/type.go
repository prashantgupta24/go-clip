package systray

import (
	"sync"

	"github.com/getlantern/systray"
)

type clipboard struct {
	menuItemArray     []*systray.MenuItem
	nextMenuItemIndex int
	menuItemToVal     map[*systray.MenuItem]string
	valExistsMap      map[string]bool
	activeSlots       int
	mutex             sync.RWMutex
}
