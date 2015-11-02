package main

import (
	"log"

	"github.com/azr/airship.go"
)

const AppKey = "YOUR_APP_KEY"
const AppMasterSecret = "YOUR_APP_MASTER_SECRET"
const IOSDeviceToken = "YOUR_DEVICE_TOKEN"

func main() {
	app := airship.App{Key: AppKey, MasterSecret: AppMasterSecret}
	data := airship.PushData{
		Audience: airship.Audience{
			IOS: IOSDeviceToken,
		},
		Notification: airship.Notification{
			Alert: "Yo man !",
			IOS: &airship.IOS{
				Alert: "Yo man !",
				Badge: "+1",
			},
		},
		DeviceTypes: "all",
	}
	err := app.Push(data)
	if err != nil {
		log.Print(err)
	}
}
