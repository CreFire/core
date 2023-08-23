package otherStruct

import (
	"demo/tools/core"
	"encoding/json"
)

// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hashmap implements a map backed by a hash table.
//
// Elements are unordered in the map.
//
// Structure is not thread safe.
//
// Reference: http://en.wikipedia.org/wiki/Associative_array

// Map holds the elements in go's native map
type Map[K comparable, V any] struct {
	m map[K]V
}

// Clear removes all elements from the map.
func (m *Map[K, V]) Clear() {
	m.m = make(map[K]V)
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{m: make(map[K]V)}
}

// Put inserts element into the map.
func (m *Map[K, V]) Put(key K, value V) {
	m.m[key] = value
}

// Get searches the element in the map by key and returns its value or nil if key is not found in map.
// Second return parameter is true if key was found, otherwise false.
func (m *Map[K, V]) Get(key K) (value V, found bool) {
	value, found = m.m[key]
	return
}

// FromJSON populates the map from the input JSON representation.
func (m *Map[K, V]) FromJSON(data []byte) error {
	elements := make(map[K]V)
	err := json.Unmarshal(data, &elements)
	if err == nil {
		m.Clear()
		for key, value := range elements {
			m.m[key] = value
		}
	}
	return err
}
func (m *Map[K, V]) ToJSON() ([]byte, error) {
	elements := make(map[string]V)
	for key, value := range m.m {
		elements[core.ToString(key)] = value
	}
	return json.Marshal(&elements)
}
