// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package helper : log.go provides panic support.
package helper

import (
	"fmt"
	"time"
)

// WriteStdout is used to write a system message directly to stdout.
func WriteStdout(goRoutine string, namespace string, functionName string, message string) {
	fmt.Printf("%s : %s : %s : %s : %s\n", time.Now().Format("2006-01-02T15:04:05.000"), goRoutine, namespace, functionName, message)
}

// WriteStdoutf is used to write a formatted system message directly stdout.
func WriteStdoutf(goRoutine string, namespace string, functionName string, format string, a ...interface{}) {
	WriteStdout(goRoutine, namespace, functionName, fmt.Sprintf(format, a...))
}
