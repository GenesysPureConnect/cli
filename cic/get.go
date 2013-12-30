/* Copyright © 1997-2013 Interactive Intelligence, Inc.  All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"bytes"
)

var cmdGet = &Command{
	Run:   runGet,
	Usage: "get [configurationtype] [id] [optionalproperties]",
	Short: "Gets a configuration object",
	Long: `
Gets a configuration object

Examples:

  cic get user kevin.glinski extension
`,
}
var invalidProperties = []string{"securityRights", "accessRights", "administrativeRights", "mailboxProperties", "licenseProperties"}

func runGet(cmd *Command, args []string) {
	var properties string

	if len(args) < 3 {
		ValidateArgCount(2,args)
		var buffer bytes.Buffer
		objectDefaults, _ := Defaults(args[0])

		for key, _ := range objectDefaults {
			if !stringInSlice(key, invalidProperties) {
				buffer.WriteString(key + ",")
			}

		}
		properties = buffer.String()

	} else {
		ValidateArgCount(3,args)
		properties = args[2]
	}

	object, err := Get(args[0], args[1], properties)

	if err != nil {
		ErrorAndExit(err.Error())
	}

	if object != nil {
		DisplayConfigRecord(object)
	}

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
