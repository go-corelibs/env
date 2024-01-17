// Copyright (c) 2024  The Go-CoreLibs Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package env provides a collection of environment variable utilities
package env

import (
	"os"
)

var (
	_env Env = New()
)

func init() {
	_env.Import(os.Environ())
}

// Default returns a package global Env instance, populated with the existing
// os.Environ variables
func Default() (env Env) {
	env = _env
	return
}

// Len is a wrapper around the Default Env.Len
func Len() (count int) {
	count = _env.Len()
	return
}

// Clear is a wrapper around the Default Env.Clear
func Clear() {
	_env.Clear()
	return
}

// Environ is a wrapper around the Default Env.Env
func Environ() (variables []string) {
	variables = _env.Environ()
	return
}

// Clone is a wrapper around the Default Env.Clone
func Clone() (clone Env) {
	clone = _env.Clone()
	return
}

// Export is a wrapper around the Default Env.Export
func Export() (err error) {
	err = _env.Export()
	return
}

// Import is a wrapper around the Default Env.Import
func Import(environment []string) {
	_env.Import(environment)
	return
}

// Expand is a wrapper around the Default Env.Expand
func Expand(input string) (expanded string) {
	expanded = _env.Expand(input)
	return
}

// Get is a wrapper around the Default Env.Get
func Get(key string) (value string, present bool) {
	value, present = _env.Get(key)
	return
}

// Set is a wrapper around the Default Env.Set
func Set(key, value string) {
	_env.Set(key, value)
	return
}

// Bool is a wrapper around the Default Env.Bool
func Bool(key string, def bool) (state bool) {
	state = _env.Bool(key, def)
	return
}

// Int is a wrapper around the Default Env.Int
func Int(key string, def int) (number int) {
	number = _env.Int(key, def)
	return
}

// Float is a wrapper around the Default Env.Float
func Float(key string, def float64) (decimal float64) {
	decimal = _env.Float(key, def)
	return
}

// String is a wrapper around the Default Env.String
func String(key string, def string) (value string) {
	value = _env.String(key, def)
	return
}
