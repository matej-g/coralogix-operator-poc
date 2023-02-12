package controllers

func ReverseMap[K, V comparable](m map[K]V) map[V]K {
	n := make(map[V]K)
	for k, v := range m {
		n[v] = k
	}
	return n
}

func SlicesWithUniqueValuesEqual[V comparable](a, b []V) bool {
	if len(a) != len(b) {
		return false
	}

	valuesSet := make(map[V]bool, len(a))
	for _, _a := range a {
		valuesSet[_a] = true
	}

	for _, _b := range b {
		if !valuesSet[_b] {
			return false
		}
	}

	return true
}
