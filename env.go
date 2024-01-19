// Copyright (c) 2024  The Go-Curses Authors
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

package env

import (
	"maps"
	"os"
	"strconv"
	"strings"
	"sync"

	clpath "github.com/go-corelibs/path"
	rpl "github.com/go-corelibs/replace"
	"github.com/go-corelibs/slices"
	clstrings "github.com/go-corelibs/strings"
)

var _ Env = (*cEnv)(nil)

// Env is an abstraction around the standard os package env related things
// with some conveniences added like Env.Int and Env.Float
type Env interface {
	// Len returns the number of variables present
	Len() (count int)
	// Clear deletes all Env variables
	Clear()
	// Environ returns a copy of all the underlying keys and values in the
	// form of "key=value"
	Environ() (variables []string)
	// Clone returns a complete duplicate of the current Env, changes to
	// a clone do not have any effect upon the original Env
	Clone() (clone Env)
	// Export updates the actual os environment, calling os.Setenv for each
	// variable present. Export stops at the first error
	Export() (err error)
	// Import updates the Env instance with the given `environment`
	// variables, in the form of "key=value". Inputs missing the equal sign are
	// ignored and all values have any quotes trimmed. Quotations are detected
	// when the first and last characters are the same and are one of the
	// following: backtick (&96;), quote (') or double quote (")
	Import(environment []string)
	// Include applies all variables within the others to this Env instance.
	// Note that keys are not deleted and any existing keys are clobbered by
	// the others, in the order the others are given
	Include(others ...Env)
	// WriteEnvDir makes the given directory path if it doesn't exist already
	// and then for each key/value pair in the Env, creates a file named with
	// the key and the value as the contents
	WriteEnvDir(path string) (err error)
	// Expand replaces all `$key` and `${key}` references in the `input` string
	// with their corresponding `key` values. Any references not present within
	// the Env are replaced with empty strings
	Expand(input string) (expanded string)
	// Get looks for the variable `key` and if `present` returns the exact
	// `value`
	Get(key string) (value string, present bool)
	// Set updates the Env `key` with the given `value`
	Set(key, value string)
	// Unset removes the `key` from the Env
	Unset(key string)
	// Bool transforms the value associated with `key` into a boolean state. If
	// the value is not a detectable state, `def` is returned. Detection is
	// handled by the github.com/go-corelibs/strings package IsTrue and IsFalse
	// functions
	Bool(key string, def bool) (state bool)
	// Int uses strconv.Atoi to transform the value associated with `key` and
	// if not present or strconv.Atoi encountered an error, returns `def`
	Int(key string, def int) (number int)
	// Float uses strconv.ParseFloat to transform the value associated with
	// `key` and if not present or strconv.ParseFloat encountered an error,
	// returns `def`
	Float(key string, def float64) (decimal float64)
	// String uses strings.TrimSpace to transform the value associated with
	// `key` and if not present, returns `def`
	String(key string, def string) (value string)
}

// New constructs a new Env instance with no variables present
func New() (env Env) {
	env = newEnv()
	return
}

// NewImport constructs a new Env instance and calls Import with the given
// `environ` slice
func NewImport(environ []string) (env Env) {
	env = New()
	env.Import(environ)
	return
}

func newEnv() (env *cEnv) {
	env = &cEnv{
		data:  make(map[string]string),
		order: make([]string, 0),
		m:     &sync.RWMutex{},
	}
	return
}

type cEnv struct {
	data  map[string]string
	order []string
	m     *sync.RWMutex
}

func (c *cEnv) Len() (count int) {
	count = len(c.order)
	return
}

func (c *cEnv) Clear() {
	c.m.Lock()
	defer c.m.Unlock()
	c.data = make(map[string]string)
	c.order = make([]string, 0)
}

func (c *cEnv) Environ() (variables []string) {
	c.m.RLock()
	defer c.m.RUnlock()
	for _, key := range c.order {
		if value, present := c.data[key]; present {
			variables = append(variables, key+"="+value)
		}
	}
	return
}

func (c *cEnv) Clone() (clone Env) {
	c.m.RLock()
	defer c.m.RUnlock()
	cloned := &cEnv{
		data:  maps.Clone(c.data),
		order: slices.Copy(c.order),
		m:     &sync.RWMutex{},
	}
	clone = cloned
	return
}

func (c *cEnv) Export() (err error) {
	c.m.RLock()
	defer c.m.RUnlock()
	for _, key := range c.order {
		if value, present := c.data[key]; present {
			if err = os.Setenv(key, value); err == nil {
				return
			}
		}
	}
	return
}

func (c *cEnv) Import(environ []string) {
	c.m.Lock()
	defer c.m.Unlock()
	for _, input := range environ {
		if key, value, found := strings.Cut(input, "="); found && key != "" {
			c.data[key] = clstrings.TrimQuotes(value)
			if !slices.Within(key, c.order) {
				c.order = append(c.order, key)
			}
		}
	}
	return
}

func (c *cEnv) Include(others ...Env) {
	for _, other := range others {
		c.Import(other.Environ())
	}
}

func (c *cEnv) WriteEnvDir(path string) (err error) {
	c.m.RLock()
	defer c.m.RUnlock()
	if err = os.MkdirAll(path, clpath.DefaultPathPerms); err != nil {
		return
	}
	for _, key := range c.order {
		if err = os.WriteFile(path+"/"+key, []byte(c.data[key]), clpath.DefaultFilePerms); err != nil {
			return
		}
	}
	return
}

func (c *cEnv) Expand(input string) (expanded string) {
	c.m.RLock()
	defer c.m.RUnlock()
	expanded = rpl.Vars(input, c.data)
	return
}

func (c *cEnv) Get(key string) (value string, present bool) {
	c.m.RLock()
	defer c.m.RUnlock()
	value, present = c.data[key]
	return
}

func (c *cEnv) Set(key, value string) {
	c.m.Lock()
	defer c.m.Unlock()
	if key != "" {
		if _, present := c.data[key]; !present {
			c.order = append(c.order, key)
		}
		c.data[key] = value
	}
	return
}

func (c *cEnv) Unset(key string) {
	c.m.Lock()
	defer c.m.Unlock()
	if _, present := c.data[key]; present {
		delete(c.data, key)
		c.order = slices.Prune(c.order, key)
	}
	return
}

func (c *cEnv) Bool(key string, def bool) (state bool) {
	c.m.RLock()
	defer c.m.RUnlock()
	if value, present := c.Get(key); present {
		value = strings.TrimSpace(value)
		if v := clstrings.IsTrue(value); v {
			state = true
			return
		} else if v = clstrings.IsFalse(value); v {
			state = false
			return
		}
	}
	state = def
	return
}

func (c *cEnv) Int(key string, def int) (number int) {
	c.m.RLock()
	defer c.m.RUnlock()
	if value, present := c.Get(key); present {
		if v, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			number = v
			return
		}
	}
	number = def
	return
}

func (c *cEnv) Float(key string, def float64) (decimal float64) {
	c.m.RLock()
	defer c.m.RUnlock()
	if value, present := c.Get(key); present {
		if v, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			decimal = v
			return
		}
	}
	decimal = def
	return
}

func (c *cEnv) String(key string, def string) (value string) {
	c.m.RLock()
	defer c.m.RUnlock()
	if v, present := c.Get(key); present {
		value = strings.TrimSpace(v)
		return
	}
	value = def
	return
}
