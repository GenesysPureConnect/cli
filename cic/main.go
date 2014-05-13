/* Code based on project https://github.com/heroku/force
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"os"
)

var commands = []*Command{
	cmdHelp,
	cmdDefaults,
	cmdFeatures,
	cmdGet,
	cmdInteraction,
	cmdLogin,
	cmdLogout,
	cmdMakeCall,
	cmdSelect,
	cmdStatus,
	cmdVersion,
	cmdWhoami,
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		usage()
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() {
				cmd.printUsage()
			}
			if err := cmd.Flag.Parse(args[1:]); err != nil {
				os.Exit(2)
			}
			cmd.Run(cmd, cmd.Flag.Args())
			return
		}
	}
	usage()
}
