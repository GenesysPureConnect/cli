/* Code based on project https://github.com/heroku/force
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import ()

var cmdWhoami = &Command{
	Run:   runWhoami,
	Usage: "whoami ",
	Short: "Show information about currently logged in user",
	Long: `
Show information about currently logged in user

Examples:

  cic whoami
`,
}

func runWhoami(cmd *Command, args []string) {
	me, err := Whoami()

	if err != nil {
		ErrorAndExit(err.Error())
	}

	if me != nil {
		DisplayConfigRecord(me)
	}

}
