/* Copyright © 1997-2013 Interactive Intelligence, Inc.  All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */
 
package main

import ("strings"
		"bytes")

var  cmdSelect = &Command{
	Run:   runSelect,
	Usage: "select <params> FROM <configurationType> WHERE <whereClause>",
	Short: "Perform a query on a particular configuration type",
	Long: `
Examples:
	Perform a query on a particular configuration type

    cic select extension,workgroups from users where configurationId.id sw kev
	cic select from users where configurationId.id = kev
  
Valid where operators are
sw - starts with
=  - equals
ct - contains
`,
}

var operators = []string{"sw", "eq", "ct", "="}


func runSelect(cmd *Command, args []string) {
	fromIndex := getIndex("from", args)
	whereIndex := getIndex("where", args)

	configType := args[fromIndex + 1]

	var fieldsBuffer bytes.Buffer
	
	for _, value := range args[0:fromIndex] {
		fieldsBuffer.WriteString(value + ",")
	}
	selectParams := fieldsBuffer.String()
	selectParams = strings.Trim(selectParams, ",")

	var whereBuffer bytes.Buffer
	
	if(whereIndex > -1){
		whereArgs := args[whereIndex + 1 :]
		for index, value := range whereArgs {
			var separator, val string

			val = value
			
			if getIndex(value, operators) > -1 || ( (len(whereArgs) > (index + 1)) && getIndex(whereArgs[index + 1], operators) > -1){
				separator = "%20"
			}else{
				separator = ""	
			}

			if strings.ToLower(value) == "and"{
				val = ","
			}else if value == "="{
				val = "eq"
			}

			whereBuffer.WriteString(val + separator)
		}

	}
	whereClause := whereBuffer.String()
	if(len(whereClause) > 0){
		whereClause = strings.Trim(whereClause, ",")
	}

	records, err := Select(configType, selectParams, whereClause)

	if err != nil {
		ErrorAndExit(err.Error())
	}

	if records != nil {
		DisplayConfigRecords(records)
	}

}

func getIndex(word string, args[] string)(index int){

	index = 0;

	for _, key := range args {
		if(strings.ToLower(key) == word){
			return;
		}
		index++
	}

	index = -1
	return;
}
