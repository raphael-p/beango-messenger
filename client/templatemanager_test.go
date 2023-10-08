package client

import (
	"testing"

	"github.com/raphael-p/beango/test/assert"
)

func TestGetTemplate(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		template1, err := getTemplate("template-1", "<div>{{ .Name }}<div/>")
		assert.IsNil(t, err)
		assert.IsNotNil(t, template1)
	})

	t.Run("CachesByName", func(t *testing.T) {
		template1, err1 := getTemplate("template-1", "<div>{{ .Name }}<div/>")
		template1Cached, err2 := getTemplate("template-1", "<a>{ .Link }}<a/>")
		assert.IsNil(t, err1, err2)
		assert.IsNotNil(t, template1, template1Cached)
		assert.DeepEquals(t, template1, template1Cached)
	})
}
