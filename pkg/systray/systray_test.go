package systray

import (
	"testing"

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
	assert.Equal(t, clipboardInstance.activeSlots, 20)
}

func (suite *ClipTestSuite) TestAddSlots() {
	t := suite.T()
	assert.Len(t, clipboardInstance.menuItemArray, 0)
	addSlots(20, clipboardInstance)
	assert.Len(t, clipboardInstance.menuItemArray, 20)
	assert.Equal(t, getActiveSlots(clipboardInstance), 20)
}

func (suite *ClipTestSuite) TestChangeSlots() {
	t := suite.T()
	assert.Len(t, clipboardInstance.menuItemArray, 0)
	addSlots(20, clipboardInstance)
	assert.Len(t, clipboardInstance.menuItemArray, 20)

	changeActiveSlots(5, clipboardInstance)
	assert.Len(t, clipboardInstance.menuItemArray, 20)
	assert.Equal(t, 5, getActiveSlots(clipboardInstance))

	changeActiveSlots(15, clipboardInstance)
	assert.Equal(t, 15, getActiveSlots(clipboardInstance))

	changeActiveSlots(1, clipboardInstance)
	assert.Equal(t, 1, getActiveSlots(clipboardInstance))
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
