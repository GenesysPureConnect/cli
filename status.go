/* Copyright © 1997-2013 Interactive Intelligence, Inc.  All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import ()

var cmdStatus = &Command{
	Run:   runStatus,
	Usage: "status <userid> [optionalstatuskey]",
	Short: "Get or sets a user's status ",
	Long: `
Get or sets a user's status.  If the status key is specified, it will set the user's status to that key'

Examples:

  cic status kevin.glinski
  cic status kevin.glinski available
`,
}

func runStatus(cmd *Command, args []string) {

	if len(args) > 1 {
		err := SetStatus(args[0], args[1])

		if err != nil {
			ErrorAndExit(err.Error())
		}
	} else {
		ValidateArgCount(1, args)

		status, err := GetStatus(args[0])

		if err != nil {
			ErrorAndExit(err.Error())
		}

		if status != nil {
			DisplayConfigRecord(status)
		}
	}

}
