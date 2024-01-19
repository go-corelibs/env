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
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/go-corelibs/path"
)

func TestEnviron(t *testing.T) {
	Convey("New Env", t, func() {
		env := New()
		So(env, ShouldNotBeNil)
		So(env.Len(), ShouldEqual, 0)
	})

	Convey("Env.Clear", t, func() {
		env := New()
		So(env, ShouldNotBeNil)
		env.Set("key", "value")
		So(env.Len(), ShouldEqual, 1)
		env.Clear()
		So(env.Len(), ShouldEqual, 0)
	})

	Convey("Env.Environ", t, func() {
		env := New()
		So(env, ShouldNotBeNil)
		env.Set("key", "value")
		env.Set("other", "valued")
		So(env.Len(), ShouldEqual, 2)
		So(env.Environ(), ShouldEqual, []string{
			"key=value",
			"other=valued",
		})
	})

	Convey("Env.Clone", t, func() {
		env := New()
		So(env, ShouldNotBeNil)
		env.Set("key", "value")
		env.Set("other", "valued")
		So(env.Len(), ShouldEqual, 2)
		cloned := env.Clone()
		So(cloned.Len(), ShouldEqual, 2)
		So(cloned.Environ(), ShouldEqual, []string{
			"key=value",
			"other=valued",
		})
	})

	Convey("Env.Export", t, func() {
		env := newEnv()
		So(env, ShouldNotBeNil)
		env.Set("__coreutils_env_test__", "value")
		So(env.Len(), ShouldEqual, 1)
		So(env.Export(), ShouldBeNil)
		So(os.Getenv("__coreutils_env_test__"), ShouldEqual, "value")
		_ = os.Unsetenv("__coreutils_env_test__")
		env.Clear()
		env.data[""] = "value"
		env.order = append(env.order, "")
		So(env.Len(), ShouldEqual, 1)
		So(env.Export(), ShouldNotBeNil)
	})

	Convey("Env.Import", t, func() {
		env := newEnv()
		So(env, ShouldNotBeNil)
		env.Import([]string{
			"two=thing", "one=thing",
		})
		So(env.Len(), ShouldEqual, 2)
		So(env.Environ(), ShouldEqual, []string{
			"two=thing", "one=thing",
		})
	})

	Convey("Env.Include", t, func() {
		env := newEnv()
		So(env, ShouldNotBeNil)
		env.Set("two", "thing")
		env.Set("one", "thing")
		So(env.Len(), ShouldEqual, 2)
		So(env.Environ(), ShouldEqual, []string{
			"two=thing", "one=thing",
		})
		env0 := env.Clone()
		env1 := newEnv()
		env1.Set("one", "overwrite")
		env0.Include(env1)
		So(env0.Environ(), ShouldEqual, []string{
			"two=thing", "one=overwrite",
		})
		env2 := newEnv()
		env2.Set("one", "more")
		env2.Set("another", "one")
		env0 = env.Clone()
		env0.Include(env1, env2)
		So(env0.Environ(), ShouldEqual, []string{
			"two=thing", "one=more", "another=one",
		})
	})

	Convey("WriteEnvDir", t, func() {
		env := newEnv()
		So(env, ShouldNotBeNil)
		env.Set("two", "thing")
		env.Set("one", "thing")
		So(env.Len(), ShouldEqual, 2)
		So(env.Environ(), ShouldEqual, []string{
			"two=thing", "one=thing",
		})
		tempDir, err := os.MkdirTemp("", "corelibs-path.*.d")
		So(err, ShouldBeNil)
		So(tempDir, ShouldNotEqual, "")
		defer os.RemoveAll(tempDir)
		err = env.WriteEnvDir(tempDir)
		So(err, ShouldBeNil)
		So(path.IsFile(tempDir+"/one"), ShouldBeTrue)
		var data []byte
		data, err = os.ReadFile(tempDir + "/one")
		So(err, ShouldBeNil)
		So(string(data), ShouldEqual, "thing")
		So(path.IsFile(tempDir+"/two"), ShouldBeTrue)
		data, err = os.ReadFile(tempDir + "/two")
		So(err, ShouldBeNil)
		So(string(data), ShouldEqual, "thing")
		So(os.Mkdir(tempDir+"/fail", 0550), ShouldBeNil)
		So(env.WriteEnvDir(tempDir+"/fail/nope"), ShouldNotBeNil)
		So(os.Chmod(tempDir+"/fail", 0770), ShouldBeNil) // for cleanup
		So(os.Chmod(tempDir+"/one", 0440), ShouldBeNil)
		So(env.WriteEnvDir(tempDir), ShouldNotBeNil) // can't overwrite one file
	})

	Convey("Env.Expand", t, func() {
		env := New()
		So(env, ShouldNotBeNil)
		env.Set("key", "value")
		env.Set("other", "valued")
		So(env.Len(), ShouldEqual, 2)
		So(env.Expand("The ${key} is $other"), ShouldEqual, "The value is valued")
	})

	Convey("Env.Get", t, func() {
		env := New()
		So(env, ShouldNotBeNil)
		env.Set("key", "value")
		env.Set("other", "valued")
		So(env.Len(), ShouldEqual, 2)
		value, present := env.Get("key")
		So(present, ShouldBeTrue)
		So(value, ShouldEqual, "value")
		value, present = env.Get("not-a-thing")
		So(present, ShouldBeFalse)
		So(value, ShouldEqual, "")
	})

	Convey("Env.Set", t, func() {
		env := New()
		So(env, ShouldNotBeNil)
		env.Set("key", "value")
		env.Set("other", "valued")
		So(env.Len(), ShouldEqual, 2)
		So(env.Environ(), ShouldEqual, []string{
			"key=value",
			"other=valued",
		})
		env.Set("key", "thing")
		So(env.Len(), ShouldEqual, 2)
		So(env.Environ(), ShouldEqual, []string{
			"key=thing",
			"other=valued",
		})
	})

	Convey("Env.Unset", t, func() {
		env := New()
		So(env, ShouldNotBeNil)
		env.Set("key", "value")
		env.Set("other", "valued")
		So(env.Len(), ShouldEqual, 2)
		So(env.Environ(), ShouldEqual, []string{
			"key=value",
			"other=valued",
		})
		env.Unset("other")
		So(env.Len(), ShouldEqual, 1)
		So(env.Environ(), ShouldEqual, []string{
			"key=value",
		})
	})

	Convey("Env.Bool", t, func() {
		env := New()
		So(env, ShouldNotBeNil)
		env.Set("key", "1")
		env.Set("other", "valued")
		env.Set("not", "f")
		So(env.Len(), ShouldEqual, 3)
		So(env.Environ(), ShouldEqual, []string{
			"key=1",
			"other=valued",
			"not=f",
		})
		So(env.Bool("key", false), ShouldEqual, true)
		So(env.Bool("other", false), ShouldEqual, false)
		So(env.Bool("not", true), ShouldEqual, false)
	})

	Convey("Env.Int", t, func() {
		env := New()
		So(env, ShouldNotBeNil)
		env.Set("key", "1")
		env.Set("other", "valued")
		So(env.Len(), ShouldEqual, 2)
		So(env.Environ(), ShouldEqual, []string{
			"key=1",
			"other=valued",
		})
		So(env.Int("key", 10), ShouldEqual, 1)
		So(env.Int("other", 10), ShouldEqual, 10)
	})

	Convey("Env.Float", t, func() {
		env := New()
		So(env, ShouldNotBeNil)
		env.Set("key", "1.0")
		env.Set("other", "valued")
		So(env.Len(), ShouldEqual, 2)
		So(env.Environ(), ShouldEqual, []string{
			"key=1.0",
			"other=valued",
		})
		So(env.Float("key", 1.1), ShouldEqual, 1.0)
		So(env.Float("other", 1.1), ShouldEqual, 1.1)
	})

	Convey("Env.String", t, func() {
		env := New()
		So(env, ShouldNotBeNil)
		env.Set("key", "value")
		So(env.Len(), ShouldEqual, 1)
		So(env.Environ(), ShouldEqual, []string{
			"key=value",
		})
		So(env.String("key", "one"), ShouldEqual, "value")
		So(env.String("other", "one"), ShouldEqual, "one")
	})
}
