package validate

import (
	"testing"

	"github.com/raphael-p/beango/test/assert"
)

func TestUniqueList(t *testing.T) {
	t.Run("NoDuplicates", func(t *testing.T) {
		list := []int{1, 2, 3, 4, 5}
		isUnique := UniqueList(list)
		assert.Equals(t, isUnique, true)
	})

	t.Run("Duplicates", func(t *testing.T) {
		list := []string{"foo", "bar", "baz", "foo"}
		isUnique := UniqueList(list)
		assert.Equals(t, isUnique, false)
	})

	t.Run("Empty", func(t *testing.T) {
		list := []float64{}
		isUnique := UniqueList(list)
		assert.Equals(t, isUnique, true)
	})
}
