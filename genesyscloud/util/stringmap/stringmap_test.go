package stringmap

import (
	"reflect"
	"testing"
)

func TestUnitMergeMaps(t *testing.T) {
	m1 := map[string][]int{
		"key1": {1, 2, 3},
		"key2": {4, 5, 6},
	}

	m2 := map[string][]int{
		"key3": {7, 8, 9},
		"key4": {10, 11, 12},
	}

	result := MergeMaps(m1, m2)

	// Check if the result contains all the keys from m1 and m2
	for key, value := range m1 {
		if !reflect.DeepEqual(result[key], value) {
			t.Errorf("Unexpected value for key %s in result. Expected: %v, Got: %v", key, value, result[key])
		}
	}

	for key, value := range m2 {
		if !reflect.DeepEqual(result[key], value) {
			t.Errorf("Unexpected value for key %s in result. Expected: %v, Got: %v", key, value, result[key])
		}
	}
}