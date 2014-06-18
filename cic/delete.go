/* Copyright © 1997-2013 Interactive Intelligence, Inc.  All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"strings"
    "fmt"
)

var cmdDelete = &Command{
	Run:   runDelete,
	Usage: "delete <id> FROM <configurationType>",
	Short: "Perform a delete on a particular configuration item",
	Long: `
Examples:
	Perform a query on a particular configuration type

    cic delete kevin.glinski from Users
`,
}


func Delete(objectType, id string) ( err error){
	if !strings.HasSuffix(objectType, "s") {
		objectType += "s"
	}

	icws := getIcws()
    err = icws.DeleteConfigurationRecord(objectType,id)
	return err;
}


func runDelete(cmd *Command, args []string) {
	fromIndex := getIndex("from", args)

    id := args[1];
	configType := args[fromIndex+1]

	err := Delete(configType, id)

	if err != nil {
		ErrorAndExit(err.Error())
	}

    fmt.Println(fmt.Sprintf("%s deleted\n", id))
}
