/* Copyright © 1997-2013 Interactive Intelligence, Inc.  All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import ("fmt")

var cmdVersion = &Command{
	Run:   runVersion,
	Usage: "version",
	Short: "Gets version information about the cic server ",
	Long: `Gets version information about the cic server
`,
}

func runVersion(cmd *Command, args []string) {

	version, err := GetVersion()

	if err != nil {
		ErrorAndExit(err.Error())
	}

	fmt.Printf("Major Version: %s", version.MajorVersion)
	fmt.Printf("Minor Version: %s", version.MinorVersion)
	fmt.Printf("SU: %s", version.Su)
	fmt.Printf("Product Id: %s", version.ProductId)
	fmt.Printf("Code Base Id: %s", version.CodeBaseId)
	fmt.Printf("Build: %s", version.Build)
	fmt.Printf("Release Display String: %s", version.ProductReleaseDisplayString)
	fmt.Printf("Patch Display String: %s", version.ProductPatchDisplayString)

}
