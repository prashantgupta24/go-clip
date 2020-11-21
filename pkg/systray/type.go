package systray

import "github.com/getlantern/systray"

type clipboard struct {
	menuItemArray     []*systray.MenuItem
	nextMenuItemIndex int
	menuItemToVal     map[*systray.MenuItem]string
	valExistsMap      map[string]bool
	activeSlots       int
}
