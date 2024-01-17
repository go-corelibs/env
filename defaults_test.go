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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDefaults(t *testing.T) {
	Convey("Basic wrapper checks", t, func() {
		So(Default(), ShouldEqual, _env)
		So(Len(), ShouldEqual, _env.Len())
		So(Environ(), ShouldEqual, _env.Environ())
		So(Clone().Environ(), ShouldEqual, _env.Environ())
		Import([]string{"coreutils_env_test=value"})
		env, _ := _env.(*cEnv)
		value, ok := env.data["coreutils_env_test"]
		So(ok, ShouldBeTrue)
		So(value, ShouldEqual, "value")
		So(Expand("${coreutils_env_test}"), ShouldEqual, "value")
		value, ok = Get("coreutils_env_test")
		So(ok, ShouldBeTrue)
		So(value, ShouldEqual, "value")
		Set("coreutils_env_test", "1")
		value, ok = Get("coreutils_env_test")
		So(ok, ShouldBeTrue)
		So(value, ShouldEqual, "1")
		So(Bool("coreutils_env_test", false), ShouldEqual, true)
		So(Int("coreutils_env_test", 10), ShouldEqual, 1)
		So(Float("coreutils_env_test", 1.1), ShouldEqual, 1.0)
		So(String("coreutils_env_test", "one"), ShouldEqual, "1")
		Clear() // do this last
		So(Len(), ShouldEqual, 0)
		So(Export(), ShouldBeNil)
	})
}
