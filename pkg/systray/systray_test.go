package systray

import (
	"strconv"
	"testing"
	"time"

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
	addSlots(30, clipboardInstance)
	assert.Len(t, clipboardInstance.menuItemArray, 30)

	changes := make(chan string, 10)
	stopCh := make(chan struct{})
	go monitorClipboard(clipboardInstance, stopCh, changes)

	for i := 0; i < 30; i++ {
		assert.Len(t, clipboardInstance.menuItemToVal, i)
		assert.Len(t, clipboardInstance.valExistsMap, i)
		changes <- "hello" + strconv.Itoa(i)
		time.Sleep(time.Millisecond * 10)
		// fmt.Println("clipboardInstance.menuItemToVal", clipboardInstance.menuItemToVal)
		assert.Len(t, clipboardInstance.menuItemToVal, i+1)
		assert.Len(t, clipboardInstance.valExistsMap, i+1)
		assert.Contains(t, clipboardInstance.valExistsMap, "hello"+strconv.Itoa(i))
	}

	// assert.Equal(t, "hello", )

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
