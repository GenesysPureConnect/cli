/* Code based on project https://github.com/heroku/force
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"os"
)

const (
	LF = 10
)

func ErrorAndExit(format string, args ...interface{}) {
	
	if len(format) > 0 && format[0] == LF {
		fmt.Fprintf(os.Stderr, format[1:]+"\n", args...)
	} else {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("ERROR: %s\n", format), args...)
	}
	os.Exit(1)
}

func ValidateArgCount(expected int, args []string) {
	if len(args) != expected {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("ERROR: Expected %d args, received %d\n", expected, len(args)))
		os.Exit(1)
	}
}
