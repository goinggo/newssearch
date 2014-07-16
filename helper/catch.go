// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package helper : catch.go provides panic support.
package helper

import (
	"fmt"
	"runtime"
)

// CatchPanic is used to catch any Panic and log exceptions to Stdout. It will also write the stack trace.
func CatchPanic(err *error, goRoutine string, namespace string, functionName string) {
	if r := recover(); r != nil {
		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		WriteStdoutf(goRoutine, namespace, functionName, "PANIC Defered [%v] : Stack Trace : %v", r, string(buf))

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}
