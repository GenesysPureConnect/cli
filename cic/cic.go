/* Code based on project https://github.com/heroku/force
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type ConfigRecord map[string]interface{}

func Defaults(configurationType string) (defaults ConfigRecord, err error) {

	server, session := getServerAndSession()

	body, err := httpGet(server + "/icws/" + session + "/configuration/defaults/" + configurationType)

	if err != nil {
		ErrorAndExit(err.Error())
	}

	err = json.Unmarshal(body, &defaults)
	return

}

func Get(configurationType, id, properties string) (defaults ConfigRecord, err error) {

	server, session := getServerAndSession()

	if !strings.HasSuffix(configurationType, "s") {
		configurationType += "s"
	}

	body, err := httpGet(server + "/icws/" + session + "/configuration/" + configurationType + "/" + id + "?select=" + properties)

	if err != nil {
		ErrorAndExit(err.Error())
	}

	err = json.Unmarshal(body, &defaults)
	return

}

func GetStatus(userId string) (defaults ConfigRecord, err error) {

	server, session := getServerAndSession()

	body, err := httpGet(server + "/icws/" + session + "/status/user-statuses/" + userId)

	if err != nil {
		ErrorAndExit(err.Error())
	}

	err = json.Unmarshal(body, &defaults)
	return

}

func SetStatus(userId, statusKey string) (err error) {

	server, session := getServerAndSession()

	var statusData = map[string]string{
		"statusId": statusKey,
	}
	_, err = httpPut(server+"/icws/"+session+"/status/user-statuses/"+userId, statusData)

	if err != nil {
		ErrorAndExit(err.Error())
	}

	return

}

func GetVersion() (version ConfigRecord, err error) {

	server, err := Config.Load("current", "server")
	if err != nil {
		ErrorAndExit(err.Error())
	}

	body, err := httpGet(server + "/icws/connection/version")

	if err != nil {
		ErrorAndExit(err.Error())
	}

	err = json.Unmarshal(body, &version)

	return

}

func GetFeatures() (features []string, err error) {

	server, err := Config.Load("current", "server")
	if err != nil {
		ErrorAndExit(err.Error())
	}

	body, err := httpGet(server + "/icws/connection/features")

	if err != nil {
		ErrorAndExit(err.Error())
	}

	var featureMap map[string][]map[string]interface{}
	err = json.Unmarshal(body, &featureMap)

	features = make([]string, len(featureMap["featureInfoList"]))

	i := 0
	for _, value := range featureMap["featureInfoList"] {
		features[i] = fmt.Sprintf("%v", value["featureId"])
		i++
	}

	return

}

func Select(objectType, selectFields, where string) (records []ConfigRecord, err error) {

	if !strings.HasSuffix(objectType, "s") {
		objectType += "s"
	}

	server, session := getServerAndSession()

	var selectString string
	if selectFields == "*" {
		selectString = ""
	} else {
		selectString = "select=" + selectFields
	}

	var whereString string
	if len(where) == 0 {
		whereString = ""
	} else {
		whereString = "&where=" + where
	}

	body, err := httpGet(server + "/icws/" + session + "/configuration/" + objectType + "?" + selectString + whereString)
	if err != nil {
		return
	}

	var result map[string][]ConfigRecord
	err = json.Unmarshal(body, &result)

	records = result["items"]
	return

}

func Whoami() (me ConfigRecord, err error) {

	server, session := getServerAndSession()

	username, err := Config.Load("current", "username")
	if err != nil {
		return
	}

	body, err := httpGet(server + "/icws/" + session + "/configuration/users/" + username + "?select=extension,defaultWorkstation,statusText,roles,skills,workgroups")

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &me)
	return

}

func MakeCall(target string) (result ConfigRecord, err error) {

	server, session := getServerAndSession()

	var callData = map[string]string{
		"__type": "urn:inin.com:interactions:createCallParameters",
		"target": target,
	}
	body, err, _ := httpPost(server+"/icws/"+session+"/interactions", callData)

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &result)
	return

}

func InteractionAction(action, interactionId, attribute string) (result ConfigRecord, err error) {

	server, session := getServerAndSession()
	var body []byte
	if action == "get" {
		body, err = httpGet(server + "/icws/" + session + "/interactions/" + interactionId + "?select=" + attribute)

		if err != nil {
			return
		}

		err = json.Unmarshal(body, &result)

	} else if action == "set" {
		/*body, err, _ := httpPost(server + "/icws/" + session + "/interactions/" + interactionId " , callData)

		  err = json.Unmarshal(body, &result)
		*/
	} else {
		var isOn = "false"
		if attribute == "on" || attribute == "yes" || attribute == "1" {
			isOn = "true"
		}
		var callData = map[string]string{
			"on": isOn,
		}

		if len(attribute) == 0 {
			callData = nil
		}

		_, err, _ = httpPost(server+"/icws/"+session+"/interactions/"+interactionId+"/"+action, callData)

	}

	return

}

func Login(server, username, password string) (token string, session string, cookie string, err error) {

	var loginData = map[string]string{
		"__type":          "urn:inin.com:connection:icAuthConnectionRequestSettings",
		"applicationName": "CLI",
		"userID":          username,
		"password":        password,
	}
	body, err, cookie := httpPost(server+"/icws/connection", loginData)

	if err == nil {
		var returnData map[string]string
		json.Unmarshal(body, &returnData)
		token = returnData["csrfToken"]
		session = returnData["sessionId"]
	} else {
		fmt.Println(fmt.Sprintf("ERROR: %s\n", err.Error()))
	}
	return
}

func getServerAndSession() (server, session string) {
	var err error
	server, err = Config.Load("current", "server")
	if err != nil {
		ErrorAndExit(err.Error())
	}

	session, err = Config.Load("current", "session")
	if err != nil {
		ErrorAndExit(err.Error())
	}
	return
}

func httpGet(url string) (body []byte, err error) {

	req, err := httpRequest("GET", url, nil)

	if err != nil {
		return
	}

	response, err := httpClient().Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode == 401 {
		err = errors.New("authorization expired, please run `cic login`")
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	if response.StatusCode/100 != 2 {
		err = errors.New(createErrorMessage(response.StatusCode, body))

		return
	}

	return
}

func createErrorMessage(statusCode int, body []byte) string {

	var errorDescription string

	switch statusCode {
	case 400:
		errorDescription = "Bad Request (400)"
	case 401:
		errorDescription = "Unauthorized (401)"
	case 403:
		errorDescription = "Forbidden (403)"
	case 404:
		errorDescription = "Not Found (404)"
	case 410:
		errorDescription = "Gone (410)"
	case 500:
		errorDescription = "Internal Server Error (500)"
	}

	var message map[string]interface{}
	json.Unmarshal(body, &message)
    
	return errorDescription + ": " + message["errorId"].(string) + " " + message["message"].(string)
}

func httpPost(url string, attrs map[string]string) (body []byte, err error, cookie string) {

	rbody, _ := json.Marshal(attrs)
	req, err := httpRequest("POST", url, bytes.NewReader(rbody))
	if err != nil {
		return
	}

	response, err := httpClient().Do(req)
	if err != nil {
		return

	}
	defer response.Body.Close()
	if response.StatusCode == 401 {
		err = errors.New("authorization expired, please run `cic login`")
		return
	}
	body, err = ioutil.ReadAll(response.Body)

	if response.StatusCode/100 != 2 {
		err = errors.New(createErrorMessage(response.StatusCode, body))

		return
	}
    

	if response.Header["Set-Cookie"] != nil {
		cookie = response.Header["Set-Cookie"][0]
	}
	return
}

func httpPut(url string, attrs map[string]string) (body []byte, err error) {

	rbody, _ := json.Marshal(attrs)
	req, err := httpRequest("PUT", url, bytes.NewReader(rbody))
	if err != nil {
		return
	}

	response, err := httpClient().Do(req)
	if err != nil {
		return

	}
	defer response.Body.Close()
	if response.StatusCode == 401 {
		err = errors.New("authorization expired, please run `cic login`")
		return
	}
	body, err = ioutil.ReadAll(response.Body)

	if response.StatusCode/100 != 2 {
		err = errors.New(createErrorMessage(response.StatusCode, body))

		return
	}

	return
}

func httpClient() (client *http.Client) {

	client = &http.Client{}
	return
}

func httpRequest(method, url string, body io.Reader) (request *http.Request, err error) {
	request, err = http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept-Language", "en-us")

	cookie, configerr := Config.Load("current", "cookie")
	if configerr == nil {
		request.Header.Add("Cookie", cookie)
	} else {
		return
	}

	token, configerr := Config.Load("current", "token")

	if configerr == nil {
		request.Header.Add("ININ-ICWS-CSRF-Token", token)
	} else {
		return
	}

	// request.Header.Add("User-Agent", fmt.Sprintf("cic cli (%s-%s)", runtime.GOOS, runtime.GOARCH))
	return
}
