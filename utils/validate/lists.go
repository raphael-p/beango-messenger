package validate

func UniqueList[K comparable](s []K) bool {
	seen := make(map[K]bool)
	for _, str := range s {
		if _, ok := seen[str]; !ok {
			return false
		}
		seen[str] = true
	}
	return true
}
