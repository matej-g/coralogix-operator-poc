package api

import "google.golang.org/protobuf/types/known/wrapperspb"

type Diff struct {
	Name            string
	Desired, Actual interface{}
}

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

func GetKeys[K, V comparable](m map[K]V) []K {
	result := make([]K, 0)
	for k := range m {
		result = append(result, k)
	}
	return result
}

func StringSliceToWrappedStringSlice(arr []string) []*wrapperspb.StringValue {
	result := make([]*wrapperspb.StringValue, 0, len(arr))
	for _, s := range arr {
		result = append(result, wrapperspb.String(s))
	}
	return result
}

func WrappedStringSliceToStringSlice(arr []*wrapperspb.StringValue) []string {
	result := make([]string, 0, len(arr))
	for _, s := range arr {
		result = append(result, s.Value)
	}
	return result
}
