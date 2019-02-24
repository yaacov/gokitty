// Copyright 2019 Yaacov Zamir <kobi.zamir@gmail.com>
// and other contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main
package main

// Store holds the key value pairs.
type Store struct {
	// key value store.
	vals map[string]string
}

func newStore() *Store {
	s := Store{
		vals: make(map[string]string),
	}

	return &s
}

func (s Store) list() map[string]string {
	return s.vals
}

func (s Store) get(k string) (string, bool) {
	val, ok := s.vals[k]

	return val, ok
}

func (s *Store) upsert(k string, v string) {
	s.vals[k] = v
}

func (s *Store) delete(k string) {
	delete(s.vals, k)
}
