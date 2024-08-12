package main

import (
	"fmt"
)

func asAny(a any) any { return a }

func main() {
	fmt.Println("cheese")
	var arr []int
	value := ArrGet(arr, 2, 0) // Returns 0, as the array is nil
	fmt.Println(value)

	var m map[string]int
	mapValue := MapGet(m, "key", -1) // Returns -1, as the map is nil
	fmt.Println(mapValue)

	var ptr *float64
	dereferencedValue := Deref(ptr, 0.0) // Returns 0.0, as the pointer is nil
	fmt.Println(dereferencedValue)

	var iface interface{} = "hello"
	unfacedValue := Unface[int](iface, 10) // Returns 10, as the type does not match
	fmt.Println(unfacedValue)

}

//package SafeOps

func isNil(i any) bool {
	if i == nil {
		return true
	}
	return false
}

// ArrGet retrieves the element at the specified index from a slice.
// Negative indices count from the end of the slice.
func ArrGet[T any](slice []T, index int, defaultValue ...T) T {
	var zt T
	if len(defaultValue) > 0 {
		zt = defaultValue[0]
	}
	a := &[1]T{zt}
	arrGet(a, slice, index, zt)
	return a[0]
}

func arrGet[T any](a *[1]T, slice []T, index int, defaultValue T) {
	defer func() {
		if r := recover(); r != nil {
			a[0] = defaultValue
		}
	}()
	n := len(slice)
	if index < 0 {
		index = n + index
	}
	if index < 0 || index >= n {
		a[0] = defaultValue
	}
	a[0] = slice[index]
	if isNil(a[0]) {
		a[0] = defaultValue
	}
}

func ArrGetSlice[T any](slice [][]T, index int) []T {
	return ArrGetFrom(slice, index, func() []T { return []T{} })
}

func ArrGetMap[K comparable, V any](m []map[K]V, index int) map[K]V {
	return ArrGetFrom(m, index, func() map[K]V { return make(map[K]V) })
}

func ArrGetFrom[T any](slice []T, index int, defaultFn func() T) T {
	var zt T
	a := &[1]T{zt}
	arrGetFrom(a, slice, index, defaultFn)
	return a[0]
}

func arrGetFrom[T any](a *[1]T, slice []T, index int, defaultFn func() T) {
	defer func() {
		if r := recover(); r != nil {
			a[0] = defaultFn()
		}
	}()
	n := len(slice)
	if index < 0 {
		index = n + index
	}
	if index < 0 || index >= n {
		a[0] = defaultFn()
		return
	}
	a[0] = slice[index]
	if isNil(a[0]) {
		a[0] = defaultFn()
	}
}

// ArrSet sets the value at the specified index in a slice.
// Negative indices count from the end of the slice.
// If the index is out of range, it appends to the end of the slice.
func ArrSet[T any](slice []T, index int, value T) []T {
	s := &[1][]T{slice}
	arrSet(s, slice, index, value)
	return s[0]
}
func arrSet[T any](s *[1][]T, slice []T, index int, value T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	n := len(slice)
	if index < 0 {
		index = n + index
	}
	if index < 0 {
		s[0] = slice
		return
	}
	if index >= n {
		temp := make([]T, index+1)
		copy(temp, slice)
		temp[index] = value
		s[0] = temp
		return
	}
	slice[index] = value
	s[0] = slice
	return
}

// MapGet retrieves a value from the map using the key, or returns defaultValue if the key does not exist or the map is nil.
func MapGet[K comparable, V any](m map[K]V, key K, defaultValue ...V) V {
	var zt V
	if len(defaultValue) > 0 {
		zt = defaultValue[0]
	}
	v := &[1]V{zt}
	mapGet(v, m, key, zt)
	return v[0]
}
func mapGet[K comparable, V any](v *[1]V, m map[K]V, key K, defaultValue V) {
	defer func() {
		if r := recover(); r != nil {
			v[0] = defaultValue
		}
	}()
	if m == nil {
		v[0] = defaultValue
		return
	}
	value, ok := m[key]
	if !ok {
		v[0] = defaultValue
		return
	}
	if isNil(v[0]) {
		v[0] = defaultValue
		return
	}
	v[0] = value
}

func MapGetSlice[T any,K comparable](m map[K][]T, key K) []T {
	return MapGetFrom(m, key, func()[]T{ return []T{}})
}

func MapGetMap[Q comparable,K comparable, V any](m map[Q]map[K]V, key Q) map[K]V {
	return MapGetFrom(m, key, func() map[K]V { return make(map[K]V) })
}

func MapGetFrom[K comparable, V any](m map[K]V, key K, defaultFn func() V) V {
	var zt V
	a := &[1]V{zt}
	mapGetFrom(a, m, key, defaultFn)
	return a[0]
}

func mapGetFrom[K comparable, V any](a *[1]V, m map[K]V, key K, defaultFn func() V) {
	defer func() {
		if r := recover(); r != nil {
			a[0] = defaultFn()
		}
	}()
	a[0] = m[key]
	if isNil(a[0]) {
		a[0] = defaultFn()
	}
}

// Deref dereferences a pointer to get the value it points to, or returns defaultValue if the pointer is nil.
func Deref[T any](ptr *T, defaultValue ...T) T {
	var zt T
	if len(defaultValue) > 0 {
		zt = defaultValue[0]
	}
	t := &[1]T{zt}
	deref(t, ptr, zt)
	return t[0]
}
func deref[T any](t *[1]T, ptr *T, defaultValue T) {
	defer func() {
		if r := recover(); r != nil {
			t[0] = defaultValue
		}
	}()
	if ptr == nil {
		t[0] = defaultValue
		return
	}
	t[0] = *ptr
	if isNil(t[0]) {
		t[0] = defaultValue
	}
}

// Unface retrieves a value from an interface based on its type,
// or returns the user-provided defaultValue if the interface is nil or not of the expected type.
func Unface[T any](val interface{}, defaultValue ...T) T {
	var zt T
	if len(defaultValue) > 0 {
		zt = defaultValue[0]
	}
	t := &[1]T{zt}
	unface(t, val, zt)
	return t[0]
}
func unface[T any](t *[1]T, val interface{}, defaultValue T) {
	defer func() {
		if r := recover(); r != nil {
			t[0] = defaultValue
		}
	}()
	if val == nil {
		t[0] = defaultValue
		return
	}
	switch v := val.(type) {
	case T:
		t[0] = v
	default:
		t[0] = defaultValue
	}
}
