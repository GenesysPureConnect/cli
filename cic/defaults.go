/* Copyright © 1997-2013 Interactive Intelligence, Inc.  All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */
 
package main

import ()

var cmdDefaults = &Command{
	Run:   runDefaults,
	Usage: "defaults [configurationtype] ",
	Short: "Show the default values for a new item of the configuration type",
	Long: `
Show the default values for a new item of the configuration type

Examples:

  cic defaults user
`,
}

func runDefaults(cmd *Command, args []string) {
	object, err := Defaults(args[0])

	if err != nil {
		ErrorAndExit(err.Error())
	}

	if object != nil {
		DisplayConfigRecord(object)
	}

}
