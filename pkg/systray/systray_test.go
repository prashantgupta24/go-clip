package systray

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/getlantern/systray"
	"github.com/prashantgupta24/go-clip/clip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ClipTestSuite struct {
	suite.Suite
}

func (suite *ClipTestSuite) SetupTest() {
	clipboardInstance = nil
	initInstance()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ClipTestSuite))
}

func (suite *ClipTestSuite) TestInit() {
	t := suite.T()
	assert.NotNil(t, clipboardInstance)
	initInstance()
	assert.Equal(t, clipboardInstance.activeSlots, 0)
}

func (suite *ClipTestSuite) TestInitializeClipBoard() {
	t := suite.T()
	initializeClipBoard()
	assert.Equal(t, 10, clipboardInstance.activeSlots)
}

func (suite *ClipTestSuite) TestAddSlots() {
	t := suite.T()
	assert.Len(t, clipboardInstance.menuItemArray, 0)
	addSlots(20, clipboardInstance)
	assert.Len(t, clipboardInstance.menuItemArray, 20)
	assert.Equal(t, getActiveSlots(clipboardInstance), 20)
	assert.Equal(t, 20, clipboardInstance.activeSlots)

	addSlots(20, clipboardInstance)
	assert.Len(t, clipboardInstance.menuItemArray, 40)
	assert.Equal(t, getActiveSlots(clipboardInstance), 40)
	assert.Equal(t, 40, clipboardInstance.activeSlots)

}

func (suite *ClipTestSuite) TestChangeSlots() {
	t := suite.T()
	assert.Len(t, clipboardInstance.menuItemArray, 0)
	addSlots(20, clipboardInstance)
	assert.Len(t, clipboardInstance.menuItemArray, 20)
	assert.Equal(t, 20, clipboardInstance.activeSlots)

	for i := 0; i < 20; i++ {
		changetTo := rand.Intn(20) + 1
		changeActiveSlots(changetTo, clipboardInstance)
		assert.Equal(t, changetTo, getActiveSlots(clipboardInstance))
	}
}

func (suite *ClipTestSuite) TestObfuscateVal() {
	t := suite.T()
	menuItem := menuItem{
		instance:     &systray.MenuItem{},
		subMenuItems: make(map[subMenu]*systray.MenuItem),
	}
	menuItem.subMenuItems[obfuscateMenu] = &systray.MenuItem{}

	//test1
	clipboardInstance.menuItemToVal[menuItem.instance] = "test_value"
	obfuscateVal(clipboardInstance, menuItem)

	assert.Equal(t, "test******", getTitle(menuItem))
	assert.Equal(t, "test******", getToolTip(menuItem))
	assert.True(t, menuItem.subMenuItems[obfuscateMenu].Checked())

	//test2
	clipboardInstance.menuItemToVal[menuItem.instance] = "test"
	obfuscateVal(clipboardInstance, menuItem)

	assert.Equal(t, "test", getTitle(menuItem))
	assert.Equal(t, "test", getToolTip(menuItem))
	assert.True(t, menuItem.subMenuItems[obfuscateMenu].Checked())

	//test3
	clipboardInstance.menuItemToVal[menuItem.instance] = "this_is_a_big_value"
	obfuscateVal(clipboardInstance, menuItem)

	assert.Equal(t, "this***************", getTitle(menuItem))
	assert.Equal(t, "this***************", getToolTip(menuItem))
	assert.True(t, menuItem.subMenuItems[obfuscateMenu].Checked())
}

func (suite *ClipTestSuite) TestAcceptVal() {
	t := suite.T()
	addSlots(20, clipboardInstance)
	for i := 0; i < 20; i++ {
		menuItem := clipboardInstance.menuItemArray[i]
		val1 := "test_message_" + strconv.Itoa(i)

		assert.Equal(t, "", getTitle(menuItem))
		assert.Equal(t, "(empty slot)", getToolTip(menuItem))
		assert.False(t, clipboardInstance.valExistsMap[val1])
		assert.NotContains(t, clipboardInstance.menuItemToVal, menuItem.instance)
		assert.True(t, menuItem.subMenuItems[obfuscateMenu].Disabled())
		assert.True(t, menuItem.subMenuItems[pinMenu].Disabled())

		acceptVal(clipboardInstance, menuItem, val1)

		assert.Equal(t, val1, getTitle(menuItem))
		assert.Equal(t, val1, getToolTip(menuItem))
		assert.True(t, clipboardInstance.valExistsMap[val1])
		assert.Contains(t, clipboardInstance.menuItemToVal, menuItem.instance)
		assert.Equal(t, val1, clipboardInstance.menuItemToVal[menuItem.instance])
		assert.False(t, menuItem.subMenuItems[obfuscateMenu].Disabled())
		assert.False(t, menuItem.subMenuItems[pinMenu].Disabled())
	}
}

func (suite *ClipTestSuite) TestSubstituteMenuItem() {
	t := suite.T()
	addSlots(20, clipboardInstance)
	menuItem := menuItem{
		instance:     &systray.MenuItem{},
		subMenuItems: make(map[subMenu]*systray.MenuItem),
	}
	menuItem.subMenuItems[pinMenu] = &systray.MenuItem{}
	menuItem.subMenuItems[obfuscateMenu] = &systray.MenuItem{}
	menuItem.subMenuItems[obfuscateMenu].Check()

	existingMenuItem := getExistingSlotToReplace()

	valNew := "test_value_new"
	valExisting := "test_value_existing"

	acceptVal(clipboardInstance, menuItem, valNew)
	acceptVal(clipboardInstance, existingMenuItem, valExisting)

	assert.NotNil(t, existingMenuItem)
	assert.False(t, existingMenuItem.instance.Checked())
	assert.False(t, menuItem.instance.Checked())
	assert.Equal(t, valNew, getTitle(menuItem))
	assert.Equal(t, valNew, getToolTip(menuItem))
	assert.Equal(t, valExisting, getToolTip(existingMenuItem))
	assert.Equal(t, valExisting, getToolTip(existingMenuItem))

	substituteMenuItem(clipboardInstance, menuItem)

	assert.True(t, existingMenuItem.instance.Checked())
	assert.False(t, menuItem.instance.Checked())
	assert.Equal(t, valExisting, getTitle(menuItem))
	assert.Equal(t, valExisting, getToolTip(menuItem))
	assert.Equal(t, "test**********", getTitle(existingMenuItem))
	assert.Equal(t, "test**********", getToolTip(existingMenuItem))

}
func (suite *ClipTestSuite) TestSlotChannels() {
	t := suite.T()
	addSlots(20, clipboardInstance)

	for i := 0; i < clipboardInstance.activeSlots; i++ {

		menuItem := clipboardInstance.menuItemArray[i]
		obfuscateMenu := menuItem.subMenuItems[obfuscateMenu]

		//obfuscate
		obfuscateMenu.ClickedCh <- struct{}{}
		time.Sleep(time.Millisecond * 10)
		clipboardInstance.mutex.RLock()
		assert.True(t, obfuscateMenu.Checked())
		clipboardInstance.mutex.RUnlock()

		//unobfuscate
		obfuscateMenu.ClickedCh <- struct{}{}
		time.Sleep(time.Millisecond * 10)
		clipboardInstance.mutex.RLock()
		assert.False(t, obfuscateMenu.Checked())
		clipboardInstance.mutex.RUnlock()

		//pin
		existingMenuItem := getExistingSlotToReplace()
		pinMenuOrg := menuItem.subMenuItems[pinMenu]
		pinMenuOrg.ClickedCh <- struct{}{}
		time.Sleep(time.Millisecond * 10)
		clipboardInstance.mutex.RLock()
		assert.NotNil(t, existingMenuItem)
		assert.True(t, existingMenuItem.instance.Checked())
		if i != 0 { //first slot gets replaced by itself
			assert.False(t, menuItem.instance.Checked())
		}
		clipboardInstance.mutex.RUnlock()

		//unpin
		existingPinMenu := existingMenuItem.subMenuItems[pinMenu]
		existingPinMenu.ClickedCh <- struct{}{}
		time.Sleep(time.Millisecond * 10)
		clipboardInstance.mutex.RLock()
		assert.False(t, existingMenuItem.instance.Checked())
		clipboardInstance.mutex.RUnlock()

	}
}

func getExistingSlotToReplace() menuItem {
	for i := 0; i < clipboardInstance.activeSlots; i++ {
		existingMenuItem := clipboardInstance.menuItemArray[i]
		if !existingMenuItem.instance.Disabled() && !existingMenuItem.instance.Checked() {
			return existingMenuItem
		}
	}
	return menuItem{}
}
func (suite *ClipTestSuite) TestClearSlots() {
	t := suite.T()
	addSlots(20, clipboardInstance)
	clearSlots(clipboardInstance.menuItemArray)
	assert.Equal(t, 0, clipboardInstance.nextMenuItemIndex)
}

func (suite *ClipTestSuite) TestTruncateVal() {
	t := suite.T()
	val := "this_will_be_truncated_into_something_small"
	truncVal := truncateVal(clipboardInstance, val)
	assert.Equal(t, val[:clipboardInstance.truncateLength]+"... ("+strconv.Itoa(len(val))+" chars)", truncVal)
}

func (suite *ClipTestSuite) TestClickVal() {
	t := suite.T()

	changes := make(chan string, 1)
	stopCh := make(chan struct{})
	go clip.Monitor(time.Millisecond*50, stopCh, changes)

	addSlots(1, clipboardInstance)
	menuItem := clipboardInstance.menuItemArray[0].instance

	for i := 0; i < 10; i++ {
		val := "test_val_" + strconv.Itoa(i)
		clipboardInstance.menuItemToVal[menuItem] = val
		menuItem.ClickedCh <- struct{}{}
		value := <-changes
		assert.Equal(t, val, value)
	}
}

func (suite *ClipTestSuite) TestClipboard() {
	t := suite.T()

	addSlots(20, clipboardInstance)
	assert.Len(t, clipboardInstance.menuItemArray, clipboardInstance.activeSlots)

	changes := make(chan string, 10)
	stopCh := make(chan struct{})
	go monitorClipboard(clipboardInstance, stopCh, changes)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100; i++ {
		// v1 := len(clipboardInstance.menuItemToVal)
		// v2 := len(clipboardInstance.valExistsMap)
		// assert.Len(t, clipboardInstance.menuItemToVal, min(i, clipboardInstance.activeSlots))
		// assert.Len(t, clipboardInstance.valExistsMap, min(i, clipboardInstance.activeSlots))
		changes <- "hello" + strconv.Itoa(i)
		time.Sleep(time.Millisecond * 10)
		// assert.Len(t, clipboardInstance.menuItemToVal, v1+1)
		// assert.Len(t, clipboardInstance.valExistsMap, v2+1)
		clipboardInstance.mutex.RLock()
		assert.Contains(t, clipboardInstance.valExistsMap, "hello"+strconv.Itoa(i))
		clipboardInstance.mutex.RUnlock()

		if i%5 == 0 && i != 0 {
			time.Sleep(time.Millisecond * 10)
			changetTo := rand.Intn(20) + 1
			// fmt.Println("pclipboardInstance.nextMenuItemIndex : ", clipboardInstance.nextMenuItemIndex)
			// fmt.Println("changed to : ", changetTo)
			changeActiveSlots(changetTo, clipboardInstance)
			assert.Equal(t, changetTo, getActiveSlots(clipboardInstance))
			// fmt.Println("aclipboardInstance.nextMenuItemIndex : ", clipboardInstance.nextMenuItemIndex)
		}
	}
}

func getActiveSlots(clipboard *clipboard) int {

	activeSlots := 0
	for _, menuItem := range clipboardInstance.menuItemArray {
		if !menuItem.instance.Disabled() {
			activeSlots++
		}
	}
	return activeSlots
}

// uncomment for visual testing
// func TestMain(m *testing.M) {

// 	timeToSleep := time.Second //change accordingly
// 	initInstance()
// 	systray.Run(func() {
// 		systray.SetTemplateIcon(icon.Data, icon.Data)
// 		systray.SetTooltip("Clipboard")

// 		mQuitOrig := systray.AddMenuItem("Quit", "Quit the app")
// 		go func() {
// 			<-mQuitOrig.ClickedCh
// 			fmt.Println("Requesting quit")
// 			systray.Quit()
// 			fmt.Println("Finished quitting")
// 		}()

// 		// We can manipulate the systray in other goroutines
// 		configureMenu := systray.AddMenuItem("Configuration", "Configuration")
// 		slotsMenu := configureMenu.AddSubMenuItem("slotsMenu", "SubMenu Test (middle)")
// 		slots5 := slotsMenu.AddSubMenuItem("5", "5")
// 		slots10 := slotsMenu.AddSubMenuItem("10", "10")
// 		slots20 := slotsMenu.AddSubMenuItem("20", "20")
// 		clearMenu := configureMenu.AddSubMenuItem("Clear", "Clear")

// 		addSlots(20, clipboardInstance)

// 		changes := make(chan string, 10)
// 		stopCh := make(chan struct{})
// 		go monitorClipboard(clipboardInstance, stopCh, changes)

// 		go func() {
// 			for {
// 				select {
// 				case <-slots5.ClickedCh:
// 					// fmt.Println("changed to 5")
// 					changeActiveSlots(5, clipboardInstance)
// 				case <-slots10.ClickedCh:
// 					// fmt.Println("changed to 10")
// 					changeActiveSlots(10, clipboardInstance)
// 				case <-slots20.ClickedCh:
// 					// fmt.Println("changed to 20")
// 					changeActiveSlots(20, clipboardInstance)
// 				case <-clearMenu.ClickedCh:
// 					// fmt.Println("clear")
// 					clearSlots(clipboardInstance.menuItemArray)
// 				}
// 			}
// 		}()

// 		for i := 0; i < 100; i++ {
// 			// 	// assert.Len(, clipboardInstance.menuItemToVal, i)
// 			// 	// assert.Len(t, clipboardInstance.valExistsMap, i)
// 			time.Sleep(timeToSleep)
// 			changes <- "hello" + strconv.Itoa(i)

// 			if i%5 == 0 && i != 0 {
// 				time.Sleep(time.Millisecond * 10)
// 				changetTo := rand.Intn(15) + 5
// 				log.Println("pclipboardInstance.nextMenuItemIndex : ", clipboardInstance.nextMenuItemIndex)
// 				log.Println("changed to : ", changetTo)
// 				changeActiveSlots(changetTo, clipboardInstance)
// 				// assert.Equal(t, changetTo, getActiveSlots(clipboardInstance))
// 				log.Println("aclipboardInstance.nextMenuItemIndex : ", clipboardInstance.nextMenuItemIndex)
// 				// time.Sleep(time.Millisecond * 10)
// 			}
// 		}

// 	}, func() {})
// }
