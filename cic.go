/* Code based on project https://github.com/heroku/force
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (

		"strings"
	"github.com/interactiveintelligence/icws_golib"
)

func Defaults(configurationType string) (defaults icws_golib.ConfigRecord, err error) {

	icws := getIcws()

	defaults, err = icws.Defaults(configurationType)

	if err != nil {
		ErrorAndExit(err.Error())
	}

	return

}

func Get(configurationType, id, properties string) (record icws_golib.ConfigRecord, err error) {

	icws := getIcws()

	record, err = icws.GetConfigurationRecord(configurationType, id, properties)

	if err != nil {
		ErrorAndExit(err.Error())
	}

	return

}

func GetStatus(userId string) (defaults icws_golib.ConfigRecord, err error) {

	icws := getIcws()

	defaults,err = icws.GetStatus(userId);

	if err != nil {
		ErrorAndExit(err.Error())
	}

	return

}

func SetStatus(userId, statusKey string) (err error) {

	icws := getIcws()

	err = icws.SetStatus(userId, statusKey);

	if err != nil {
		ErrorAndExit(err.Error())
	}

	return


}

func GetVersion() (version icws_golib.ServerVersion, err error) {
	icws := getIcws()

	version, err = icws.GetVersion()

	return

}

func GetFeatures() (features []icws_golib.ServerFeature, err error) {

	icws := getIcws()

	features, err = icws.GetFeatures()

	return;
}

func Select(objectType, selectFields, where string) (records []icws_golib.ConfigRecord, err error) {

	icws := getIcws()

	records, err = icws.SelectConfigurationRecords(objectType, selectFields, where)

	if err != nil {
		ErrorAndExit(err.Error())
	}

	return

}

func Whoami() (me icws_golib.ConfigRecord, err error) {

	icws := getIcws()

	me, err = icws.GetConfigurationRecord("users", icws.UserId, "extension,defaultWorkstation,statusText,roles,skills,workgroups")

	if err != nil {
		ErrorAndExit(err.Error())
	}

	return
}

func MakeCall(target string) (result icws_golib.ConfigRecord, err error) {

	icws := getIcws()
	result, err = icws.MakeCall(target);

	if err != nil {
		ErrorAndExit(err.Error())
	}

	return

}

func InteractionAction(action, interactionId, attribute string) (result icws_golib.ConfigRecord, err error) {

	icws := getIcws()
	result, err = icws.InteractionAction(action, interactionId, attribute);

	if err != nil {
		ErrorAndExit(err.Error())
	}

	return

}

func Login(server, username, password string) (string, string, string, string, error) {

	icws := icws_golib.NewIcws();
	err := icws.Login("CLI",server, username, password);

	if err != nil {
		ErrorAndExit(err.Error())
	}

	return icws.CurrentServer, icws.CurrentToken, icws.CurrentSession, icws.CurrentCookie, err;
}

func getIcws() (icws *icws_golib.Icws) {
	var err error
	icws = icws_golib.NewIcws();

	icws.CurrentServer, err = Config.Load("current", "server")
	if err != nil {
		ErrorAndExit(err.Error())
	}

	icws.CurrentSession, err = Config.Load("current", "session")
	if err != nil {
		ErrorAndExit(err.Error())
	}

	icws.CurrentCookie, err = Config.Load("current", "cookie")
	if err != nil {
		ErrorAndExit(err.Error())
	}

	icws.CurrentToken, err = Config.Load("current", "token")
	if err != nil {
		ErrorAndExit(err.Error())
	}

	icws.UserId, err = Config.Load("current", "userid")
	if err != nil {
		ErrorAndExit(err.Error())
	}



	return
}


func getIndex(word string, args []string) (index int) {

	index = 0

	for _, key := range args {
		if strings.ToLower(key) == word {
			return
		}
		index++
	}

	index = -1
	return
}
