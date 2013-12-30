/* Code based on project https://github.com/heroku/force
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import ("fmt")

var cmdLogin = &Command{
	Run:   runLogin,
	Usage: "login",
	Short: "Log in to a CIC server",
	Long: `
Log in to a CIC server

Examples:

  cic login server user password
  
`,
}

func runLogin(cmd *Command, args []string) {
	ValidateArgCount(3,args)
	err := LoginAndSave(args[0], args[1], args[2])
	if err != nil {
		ErrorAndExit(err.Error())
	}else{
		fmt.Println("Login Successful")
	}
}

var cmdLogout = &Command{
	Run:   runLogout,
	Usage: "logout <account>",
	Short: "Log out from CIC",
	Long: `
Log out from CIC

Examples:

  cic logout 
`,
}

func runLogout(cmd *Command, args []string) {

	Config.Delete("current", "username")
	Config.Delete("current", "session")
	Config.Delete("current", "cookie")
	Config.Delete("current", "token")
	Config.Delete("current", "server")
}

func LoginAndSave(server string, username string, password string) (err error) {

	Config.Delete("current", "username")
	Config.Delete("current", "session")
	Config.Delete("current", "cookie")
	Config.Delete("current", "token")
	Config.Delete("current", "server")

	token, session, cookie, err := Login(server, username, password)
	if err != nil {
		return
	}

	Config.Save("current", "username", username)
	Config.Save("current", "session", session)
	Config.Save("current", "cookie", cookie)
	Config.Save("current", "token", token)
	Config.Save("current", "server", server)

	return
}
