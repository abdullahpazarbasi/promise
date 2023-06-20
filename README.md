# Go Promise
[![Check Status](https://github.com/abdullahpazarbasi/promise/actions/workflows/checks.yml/badge.svg)](https://github.com/abdullahpazarbasi/promise/actions/workflows/checks.yml)
[![Test Status](https://github.com/abdullahpazarbasi/promise/actions/workflows/tests.yml/badge.svg)](https://github.com/abdullahpazarbasi/promise/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/abdullahpazarbasi/promise/v3)](https://goreportcard.com/report/github.com/abdullahpazarbasi/promise/v3)
[![Go Reference](https://pkg.go.dev/badge/github.com/abdullahpazarbasi/promise/v3.svg)](https://pkg.go.dev/github.com/abdullahpazarbasi/promise/v3)
[![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/abdullahpazarbasi/promise?display_name=tag&include_prereleases)](https://github.com/abdullahpazarbasi/promise/releases)
[![GitHub](https://img.shields.io/github/license/abdullahpazarbasi/promise)](/LICENSE)

Promise / Future library for Go

## Aim

This library provides a promise interface similar to JavaScript 's

## Features

- Generics (Type Parameters) support
- Interface similar to JS 's
- Cancellable promises
- Time-out support
- Context support
- Promises which can be committed separately
- Promise maps
- Race functionality
- Automatic panic recovery
- Fluent interface / Chain API
- [No external dependency](https://pkg.go.dev/github.com/abdullahpazarbasi/promise/v3?tab=imports)

## To-Do

> no to-do

## Installation

```shell
go get github.com/abdullahpazarbasi/promise/v3
```

## Usage

### Basic Usage

```go
package foo

import (
	"context"
	"fmt"
	"github.com/abdullahpazarbasi/promise/v3"
	"time"
)

func myParentRoutine() error {
	f1 := promise.New(func(ctx context.Context) (string, error) {
		return "OK", nil
	})
	f2 := promise.New(func(ctx context.Context) (bool, error) {
		return true, nil
	})
	fmt.Println("Parallel paths are starting :")
	p1 := f1.TimeOutLimit(500 * time.Millisecond).Commit()
	p2 := f2.TimeOutLimit(400 * time.Millisecond).Commit()

	fmt.Println("Doing something on primary parallel path ...")

	fmt.Println("Waiting for one of committed promises")
	concurrentOutput1, err1 := p1.Await()
	if err1 != nil {
		return err1 // output of async function 2 ignored
	}
	fmt.Println("Waiting for the other committed promise, may be the task is already done a long time ago")
	concurrentOutput2, err2 := p2.Await()
	if err2 != nil {
		panic(err2)
	}
	fmt.Printf("Output of async function 1: %v\n", concurrentOutput1)
	fmt.Printf("Output of async function 2: %v\n", concurrentOutput2)

	return nil
}

```

### Other Examples

[examples](/examples)

## Contributors

> [List](https://github.com/abdullahpazarbasi/promise/graphs/contributors)

## License

This library licensed under [MIT](/LICENSE) License.
