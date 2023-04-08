package collections

type Set map[string]bool

func (s Set) Put(value string) {
	s[value] = true
}

func (s Set) Remove(value string) {
	delete(s, value)
}

func (s Set) Has(value string) bool {
	_, ok := s[value]
	return ok
}

func (s Set) Values() []string {
	values := make([]string, 0, len(s))
	for k := range s {
		values = append(values, k)
	}
	return values
}
