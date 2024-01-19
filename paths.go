// Copyright (c) 2024  The Go-Enjin Authors
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
	"os"
	"path/filepath"
	"strings"

	clpath "github.com/go-corelibs/path"
	"github.com/go-corelibs/slices"
)

// PATH returns the current PATH environment
func PATH() (path string) {
	path = strings.Join(PATHS(), ":")
	return
}

// PATHS returns a slice of paths from the current PATH environment
func PATHS() (paths []string) {
	if value := String("PATH", ""); value != "" {
		paths = strings.Split(value, ":")
		for i := 0; i < len(paths); i++ {
			paths[i] = strings.TrimSpace(paths[i])
		}
	}
	return
}

// PrunePATH removes the given path from the PATH environment variable
func PrunePATH(path string) (err error) {
	path = strings.TrimSpace(path)
	paths := PATHS()
	paths = slices.Prune(paths, path)
	joined := strings.Join(paths, ":")
	Set("PATH", joined)
	err = os.Setenv("PATH", joined)
	return
}

// AppendPATH moves or adds the given path to the end of the PATH
// environment variable
func AppendPATH(path string) (err error) {
	path = strings.TrimSpace(path)
	paths := PATHS()
	paths = slices.Prune(paths, path)
	paths = append(paths, path)
	joined := strings.Join(paths, ":")
	Set("PATH", joined)
	err = os.Setenv("PATH", joined)
	return
}

// PrependPATH moves or adds the given path to the start of the PATH
// environment variable
func PrependPATH(path string) (err error) {
	path = strings.TrimSpace(path)
	paths := PATHS()
	paths = slices.Prune(paths, path)
	paths = append([]string{path}, paths...)
	joined := strings.Join(paths, ":")
	Set("PATH", joined)
	err = os.Setenv("PATH", joined)
	return
}

// TidyPATH removes all paths from the PATH environment variable that
// do not actually exist on the filesystem or cannot be resolved to
// their absolute, cleaned, paths
func TidyPATH() (err error) {
	paths := PATHS()
	var existing []string
	for _, path := range paths {
		if clpath.IsDir(path) {
			if _, ee := filepath.Abs(path); ee == nil {
				existing = append(existing, filepath.Clean(path))
			}
		}
	}
	joined := strings.Join(existing, ":")
	Set("PATH", joined)
	err = os.Setenv("PATH", joined)
	return
}
