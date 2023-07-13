package path

import (
	"testing"

	"github.com/raphael-p/beango/test/assert"
)

func TestRelativeJoin(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		absolutePath, ok := RelativeJoin("leads/to/", "myfile.txt")
		assert.Equals(t, ok, true)
		assert.Contains(t, absolutePath, "/beango-messenger/utils/path/leads/to/myfile.txt")
	})

	t.Run("ExtraSlashes", func(t *testing.T) {
		absolutePath, ok := RelativeJoin("///leads//to//", "/myfile.txt//")
		assert.Equals(t, ok, true)
		assert.Contains(t, absolutePath, "/beango-messenger/utils/path/leads/to/myfile.txt")
	})
}
