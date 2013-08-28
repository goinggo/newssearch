// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"github.com/goinggo/newssearch/helper"
	"github.com/goinggo/newssearch/search"
	"os"
	"runtime"
	"strconv"
)

// main is the entry point for the program
func main() {

	defer helper.CatchPanic(nil, "System", "main", "Main")

	// Give the server enough threads to work with
	runtime.GOMAXPROCS(runtime.NumCPU() * 3)

	helper.WriteStdout("System", "main", "Main", "Started")

	// Check the parameters
	routines, command, parameters, err := CheckParameters()

	if err != nil {

		helper.WriteStdoutf("System", "main", "Main", "ERROR - Completed : %s", err)
		return
	}

	// Perform the search
	search.Run("System", routines, command, parameters)

	helper.WriteStdout("System", "main", "Main", "Completed")
}

// CheckParameters verifies we have the proper command lines parameters and
// parses them for use
func CheckParameters() (routines int, command string, parameters []string, err error) {

	helper.WriteStdoutf("System", "main", "CheckParameters", "Started : Arguments[%d]", len(os.Args))

	// We need at least 4 arguments
	if len(os.Args) < 4 {

		helper.WriteStdout("System", "main", "CheckParameters", "ERROR : Not Enough Parameters")
		search.DisplayHelpExamples("System")

		return 0, "", nil, errors.New("Not Enough Parameters")
	}

	// Capture the number of routines to use
	routines, err = strconv.Atoi(os.Args[1])

	if err != nil {

		helper.WriteStdout("System", "main", "CheckParameters", "ERROR : Routine Parameter Incorrect Type")
		search.DisplayHelpExamples("System")

		return 0, "", nil, errors.New("Incorrect Parameters")
	}

	// Capture the remaining arguments
	numberOfParameters := len(os.Args) - 3
	parameters = make([]string, numberOfParameters)

	for parameter := 0; parameter < numberOfParameters; parameter++ {

		parameters[parameter] = os.Args[parameter+3]
	}

	helper.WriteStdout("System", "main", "CheckParameters", "Completed")

	return routines, os.Args[2], parameters, nil
}
