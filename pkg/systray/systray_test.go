package systray

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/getlantern/systray"
	"github.com/go-clip/icon"
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

	changeActiveSlots(5, clipboardInstance)
	assert.Len(t, clipboardInstance.menuItemArray, 20)
	assert.Equal(t, 5, getActiveSlots(clipboardInstance))

	changeActiveSlots(15, clipboardInstance)
	assert.Equal(t, 15, getActiveSlots(clipboardInstance))

	changeActiveSlots(1, clipboardInstance)
	assert.Equal(t, 1, getActiveSlots(clipboardInstance))
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
		assert.Contains(t, clipboardInstance.valExistsMap, "hello"+strconv.Itoa(i))

		if i%5 == 0 && i != 0 {
			changetTo := rand.Intn(20) + 1
			fmt.Println("pclipboardInstance.nextMenuItemIndex : ", clipboardInstance.nextMenuItemIndex)
			fmt.Println("changed to : ", changetTo)
			changeActiveSlots(changetTo, clipboardInstance)
			assert.Equal(t, changetTo, getActiveSlots(clipboardInstance))
			fmt.Println("aclipboardInstance.nextMenuItemIndex : ", clipboardInstance.nextMenuItemIndex)
			time.Sleep(time.Millisecond * 10)
		}

		// if i == 15 {
		// 	changeActiveSlots(10, clipboardInstance)
		// 	assert.Equal(t, 0, clipboardInstance.nextMenuItemIndex)
		// }

		// if i == 50 {
		// 	changeActiveSlots(15, clipboardInstance)
		// 	assert.Equal(t, 10, clipboardInstance.nextMenuItemIndex)
		// }
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getActiveSlots(clipboard *clipboard) int {

	activeSlots := 0
	for _, menuItem := range clipboardInstance.menuItemArray {
		if !menuItem.Disabled() {
			activeSlots++
		}
	}
	return activeSlots
}

// uncomment for visual testing
func TestMain(m *testing.M) {

	timeToSleep := time.Second //change accordingly
	initInstance()
	systray.Run(func() {
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

		changes := make(chan string, 10)
		stopCh := make(chan struct{})
		go monitorClipboard(clipboardInstance, stopCh, changes)

		go func() {
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
		}()

		for i := 0; i < 100; i++ {
			// 	// assert.Len(, clipboardInstance.menuItemToVal, i)
			// 	// assert.Len(t, clipboardInstance.valExistsMap, i)
			time.Sleep(timeToSleep)
			changes <- "hello" + strconv.Itoa(i)

			if i%5 == 0 && i != 0 {
				time.Sleep(time.Millisecond * 10)
				changetTo := rand.Intn(15) + 5
				log.Println("pclipboardInstance.nextMenuItemIndex : ", clipboardInstance.nextMenuItemIndex)
				log.Println("changed to : ", changetTo)
				changeActiveSlots(changetTo, clipboardInstance)
				// assert.Equal(t, changetTo, getActiveSlots(clipboardInstance))
				log.Println("aclipboardInstance.nextMenuItemIndex : ", clipboardInstance.nextMenuItemIndex)
				// time.Sleep(time.Millisecond * 10)
			}
		}

	}, func() {})
}
