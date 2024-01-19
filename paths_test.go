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
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/go-corelibs/path"
)

func TestPaths(t *testing.T) {
	m := &sync.Mutex{}
	Convey("PATH", t, func() {
		m.Lock()
		defer m.Unlock()
		_env = NewImport(os.Environ())
		actual := os.Getenv("PATH")
		So(actual, ShouldNotEqual, "")
		So(PATH(), ShouldEqual, actual)
	})

	Convey("PATHS", t, func() {
		m.Lock()
		defer m.Unlock()
		_env = NewImport(os.Environ())
		actual := os.Getenv("PATH")
		So(actual, ShouldNotEqual, "")
		actuals := strings.Split(actual, ":")
		So(PATHS(), ShouldEqual, actuals)
	})

	Convey("PrunePATH", t, func() {
		m.Lock()
		defer m.Unlock()
		_env = NewImport(os.Environ())
		actual := os.Getenv("PATH")
		So(actual, ShouldNotEqual, "")
		actuals := strings.Split(actual, ":")
		So(len(actuals), ShouldBeGreaterThan, 0)
		So(PrunePATH(actuals[0]), ShouldBeNil)
		So(slices.Contains(PATHS(), actuals[0]), ShouldBeFalse)
	})

	Convey("AppendPATH", t, func() {
		m.Lock()
		defer m.Unlock()
		_env = NewImport(os.Environ())
		actual := os.Getenv("PATH")
		So(actual, ShouldNotEqual, "")
		actuals := strings.Split(actual, ":")
		So(len(actuals), ShouldBeGreaterThan, 0)
		pwd := path.Pwd()
		So(AppendPATH(pwd), ShouldBeNil)
		So(PATHS(), ShouldEqual, append(actuals, pwd))
		So(PrunePATH(pwd), ShouldBeNil)
	})

	Convey("PrependPATH", t, func() {
		m.Lock()
		defer m.Unlock()
		_env = NewImport(os.Environ())
		actual := os.Getenv("PATH")
		So(actual, ShouldNotEqual, "")
		actuals := strings.Split(actual, ":")
		So(len(actuals), ShouldBeGreaterThan, 0)
		pwd := path.Pwd()
		So(PrependPATH(pwd), ShouldBeNil)
		So(PATHS(), ShouldEqual, append([]string{pwd}, actuals...))
		So(PrunePATH(pwd), ShouldBeNil)
	})

	Convey("TidyPath", t, func() {
		m.Lock()
		defer m.Unlock()
		_env = NewImport(os.Environ())
		actual := os.Getenv("PATH")
		So(actual, ShouldNotEqual, "")
		actuals := strings.Split(actual, ":")
		So(len(actuals), ShouldBeGreaterThan, 0)
		pwd := path.Pwd()
		name := filepath.Base(pwd)
		rel := pwd + "/../" + name
		So(AppendPATH(rel), ShouldBeNil)
		So(PATHS(), ShouldEqual, append(actuals, rel))
		So(TidyPATH(), ShouldBeNil)
		So(PATHS(), ShouldEqual, append(actuals, pwd))
		So(PrunePATH(pwd), ShouldBeNil)
	})
}
