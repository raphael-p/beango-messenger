package path

import (
	"path/filepath"
	"runtime"
)

func RelativeJoin(elem ...string) (string, bool) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", false
	}
	path := []string{filepath.Dir(file)}
	path = append(path, elem...)
	return filepath.Join(path...), true
}
