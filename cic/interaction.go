/* Copyright © 1997-2013 Interactive Intelligence, Inc.  All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */
 
package main

import (

)

var cmdInteraction = &Command{
	Run:   runInteraction,
	Usage: "interaction <action> <interactionid>",
	Short: "Performs the action on an interaction ",
	Long: `
Performs the action on an interaction 

Examples:
    cic interaction hold 1001240023  #places the call on hold
    cic interaction set 1001240023 myAttribute=hello  #sets the call attribute myAttribute to hello
    cic interaction get 1001240023 eic_state,eic_remotename,eic_remotetn   #gets the value of the call attributes eic_state,eic_remotename,eic_remotetn

Valid Actions:
    blind-transfer
    coach
    disconnect
    hold
    join
    listen
    mute
    pause
    pickup
    make-private
    record
    send-to-voicemail

`,
}

func runInteraction(cmd *Command, args []string) {

	var attribute string

	if len(args) >= 3 {
		attribute = args[2]
	} else {
		attribute = ""
	}

	callResult, err := InteractionAction(args[0], args[1], attribute)

	if err != nil {
		ErrorAndExit(err.Error())
	}

	if callResult != nil {
		DisplayConfigRecord(callResult)
	}
}
