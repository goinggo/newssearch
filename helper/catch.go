// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package helper

import (
	"errors"
	"fmt"
	"runtime"
)

//** PUBLIC METHODS

// CatchPanic is used to catch any Panic and log exceptions to Stdout. It will also write the stack trace
//  err: A reference to the err variable to be returned to the caller. Can be nil
//  goRoutine The Go routine making the call
//  namespace: The namespace the call is being made from
//  functionName: The function makeing the call
func CatchPanic(err *error, goRoutine string, namespace string, functionName string) {
	if r := recover(); r != nil {
		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		WriteStdoutf(goRoutine, namespace, functionName, "PANIC Defered [%v] : Stack Trace : %v", r, string(buf))

		if err != nil {
			*err = errors.New(fmt.Sprintf("%v", r))
		}
	}
}
