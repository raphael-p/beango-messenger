package utils

type MapKey interface {
	string | int
}

func MapValues[K MapKey, V any](thisMap map[K]V) ([]K, []V) {
	keys := make([]K, 0, len(thisMap))
	values := make([]V, 0, len(thisMap))
	for key, value := range thisMap {
		keys = append(keys, key)
		values = append(values, value)
	}
	return keys, values
}
