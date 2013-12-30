/* Copyright © 1997-2013 Interactive Intelligence, Inc.  All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */
 
package main

import (
	"fmt"
)

var cmdMakeCall = &Command{
	Run:   runMakeCall,
	Usage: "makecall <target>",
	Short: "Places a call to the target for the currently logged in user ",
	Long: `
Examples:

  cic makecall 3172222222
`,
}

func runMakeCall(cmd *Command, args []string) {

	callResult, err := MakeCall(args[0])

	if err != nil {
		ErrorAndExit(err.Error())
	}

	if callResult != nil {
		fmt.Printf("%s\n", callResult["interactionId"])
	}
}
