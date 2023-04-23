package validate

func UniqueList[K comparable](list []K) bool {
	seen := make(map[K]bool)
	for _, str := range list {
		if _, ok := seen[str]; ok {
			return false
		}
		seen[str] = true
	}
	return true
}
