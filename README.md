[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/go-corelibs/env)
[![codecov](https://codecov.io/gh/go-corelibs/env/graph/badge.svg?token=7E99Ukp5ev)](https://codecov.io/gh/go-corelibs/env)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-corelibs/env)](https://goreportcard.com/report/github.com/go-corelibs/env)

# env - environment variable utilities

env is a collection of utilities for managing and interacting with os.Environ
related things.

# Installation

``` shell
> go get github.com/go-corelibs/env@latest
```

# Examples

## Environ

``` go
func main() {
    value := env.String("KeyName", "Default Value")
    // if "KeyName" does not exist, value == "Default Value"
    env.Set("KeyName", "1")
    number := env.Get("KeyName", 10)
    // "KeyName" exists and number == int(1)
}
```

# Go-CoreLibs

[Go-CoreLibs] is a repository of shared code between the [Go-Curses] and
[Go-Enjin] projects.

# License

```
Copyright 2024 The Go-CoreLibs Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use file except in compliance with the License.
You may obtain a copy of the license at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

[Go-CoreLibs]: https://github.com/go-corelibs
[Go-Curses]: https://github.com/go-curses
[Go-Enjin]: https://github.com/go-enjin
