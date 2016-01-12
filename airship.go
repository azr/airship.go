// Package airship allows you to easily communicate with the Urban Airship
// ( http://urbanairship.com/ ).
package airship

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var UAClient = &http.Client{}

// App represents an Urban Airship application.
type App struct {
	Key          string
	MasterSecret string
	ServerUrl    string
}

//A Push payload.
type PushData struct {
	Audience     interface{}  `json:"audience,omitempty"`
	Notification Notification `json:"notification"`
	DeviceTypes  string       `json:"device_types,omitempty"`
}
type Notification struct {
	Alert      string `json:"alert,omitempty"`      // set default alert msg
	IOS        *IOS   `json:"ios,omitempty"`        // allows to override alert msg and/or add specifics in IOS, see IOS
	Android    *Alert `json:"android,omitempty"`    // allows to override alert msg in Android
	Amazon     *Alert `json:"amazon,omitempty"`     // allows to override alert msg in Amazon
	Blackberry *Alert `json:"blackberry,omitempty"` // allows to override alert msg in Blackberry
	Mpns       *Alert `json:"mpns,omitempty"`       // allows to override alert msg in Windows Phone
	Wns        *Alert `json:"wns,omitempty"`        // allows to override alert msg in Windows
}

//Represents a simple audience setting
type Audience struct {
	IOS          string `json:"device_token,omitempty"` //the unique identifier used to target an iOS device
	Android      string `json:"apid,omitempty"`         //the unique identifier used to target an Android device
	WindowsPhone string `json:"mpns,omitempty"`         //the unique identifier used to target a Windows Phone device
	Windows      string `json:"wns,omitempty"`          //the unique identifier used to target a Windows device
	Blackberry   string `json:"device_pin,omitempty"`   //the unique identifier used to target a Blackberry device
}

//Represents a base alert message
type Alert struct {
	Alert string `json:"alert,omitempty"`
}

type IOS struct {
	Alert string `json:"alert,omitempty"`
	Sound string `json:"sound,omitempty"`
	Badge string `json:"badge,omitempty"`
}

func (app *App) deliverPayload(url string, payload io.Reader, c *http.Client) error {
	if app.ServerUrl == "" {
		app.ServerUrl = "https://go.urbanairship.com"
	}
	apiEndpoint := app.ServerUrl + url
	req, err := http.NewRequest("POST", apiEndpoint, payload)
	if err != nil {
		return err
	}
	req.SetBasicAuth(app.Key, app.MasterSecret)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.urbanairship+json; version=3;")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		respString, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return fmt.Errorf("Hit a non-200 response from UA with a status code of %d: %s\n", resp.StatusCode, respString)
	}
	return nil
}

// Takes data, marshals it, and sends it along to the broadcast API endpoint.
func (app *App) Broadcast(data PushData) error {
	return app.BroadcastWithClient(data, UAClient)
}

//Same as Broadcast but you can pass your http client.
//It's pretty convenient in appengine or when behind a secured network.
func (app *App) BroadcastWithClient(data PushData, c *http.Client) error {
	json_data, err := json.Marshal(data)
	if err != nil {
		return err
	}
	payload := bytes.NewBuffer(json_data)
	return app.deliverPayload("/api/push/broadcast", payload, c)
}

// Takes data, marshals it, and sends it along to the push API endpoint.
func (app *App) Push(data PushData) error {
	return app.PushWithClient(data, UAClient)
}

//Same as Push but you can pass your http client.
//It's pretty convenient in appengine or when behind a secured network.
func (app *App) PushWithClient(data PushData, c *http.Client) error {
	json_data, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Println("data: ", string(json_data))
	payload := bytes.NewBuffer(json_data)
	return app.deliverPayload("/api/push", payload, c)
}
