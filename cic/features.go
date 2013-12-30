/* Copyright © 1997-2013 Interactive Intelligence, Inc.  All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */
 
package main

import ()

var cmdFeatures = &Command{
	Run:   runFeatures,
	Usage: "features",
	Short: "Gets feature information about the cic server ",
	Long: `Gets feature information about the cic server
`,
}

func runFeatures(cmd *Command, args []string) {

	version, err := GetFeatures()

	if err != nil {
		ErrorAndExit(err.Error())
	}

	DisplayList(version)
}
