package validate

import (
	"testing"

	"github.com/raphael-p/beango/test/assert"
)

func TestUniqueList_NoDuplicates(t *testing.T) {
	list := []int{1, 2, 3, 4, 5}
	isUnique := UniqueList(list)
	assert.Equals(t, isUnique, true)
}

func TestUniqueList_Duplicates(t *testing.T) {
	list := []string{"foo", "bar", "baz", "foo"}
	isUnique := UniqueList(list)
	assert.Equals(t, isUnique, false)
}

func TestUniqueList_Empty(t *testing.T) {
	list := []float64{}
	isUnique := UniqueList(list)
	assert.Equals(t, isUnique, true)
}
