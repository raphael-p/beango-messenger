package server

import (
	"os"
	"testing"

	"github.com/raphael-p/beango/test/assert"
)

func TestSetup(t *testing.T) {
	t.Run("Panics", func(t *testing.T) {
		defer os.Setenv("BG_CONFIG_FILEPATH", os.Getenv("BG_CONFIG_FILEPATH"))
		_ = os.Setenv("BG_CONFIG_FILEPATH", "/not/a/real/filepath.json")
		router, ok := setup()
		assert.IsNil(t, router)
		assert.Equals(t, ok, false)
	})
}
