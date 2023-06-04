package server

import (
	"os"
	"testing"

	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/utils/logger"
)

func TestSetup(t *testing.T) {
	t.Run("Panics", func(t *testing.T) {
		defer os.Setenv("BG_CONFIG_FILEPATH", os.Getenv("BG_CONFIG_FILEPATH"))
		_ = os.Setenv("BG_CONFIG_FILEPATH", "/not/a/real/filepath.json")
		buf := logger.MockFileLogger(t)

		conn, router, ok := setup()
		assert.IsNil(t, conn)
		assert.IsNil(t, router)
		assert.Equals(t, ok, false)
		assert.Contains(t, buf.String(), "[ERROR]", "setup failed: could not open config file")
	})
}
