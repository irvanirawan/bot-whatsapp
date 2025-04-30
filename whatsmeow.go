package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/proto/waCompanionReg"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type WhatsMeowController struct {
	Client          *whatsmeow.Client
	DatabaseHandler *DatabaseHandler
}

func (meowController *WhatsMeowController) StartClient() {
	fmt.Println("Starting WhatsApp client...")

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite", "file:whatsapp.db?_foreign_keys=on&_pragma=foreign_keys(1)", dbLog)
	if err != nil {
		panic(err)
	}
	store.SetOSInfo("macOS", store.GetWAVersion())
	store.DeviceProps.PlatformType = waCompanionReg.DeviceProps_SAFARI.Enum()

	// Get the first device in the container
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	deviceStore.PushName = "RVN"
	deviceStore.Platform = "Android"
	deviceStore.Save()

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	meowController.Client = whatsmeow.NewClient(deviceStore, clientLog)
	meowController.Client.AddEventHandler(meowController.EventHandler)

	if meowController.Client.Store.ID == nil {
		// No ID stored, initiate new login
		qrChan, _ := meowController.Client.GetQRChannel(context.Background())
		err = meowController.Client.Connect()
		if err != nil {
			fmt.Println("Error connecting client:", err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Generate and store QR code
				image, _ := qrcode.Encode(evt.Code, qrcode.Medium, 256)
				base64qrcode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
				err := meowController.DatabaseHandler.UpdateQrCode(base64qrcode)
				if err != nil {
					fmt.Println("Error saving QR code:", err)
				} else {
					fmt.Println("QR code saved successfully.")
				}
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = meowController.Client.Connect()
		if err != nil {
			panic(err)
		}
	}
}

func (meowController *WhatsMeowController) EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received message:", v.Message.GetConversation())
	case *events.Connected:
		err := meowController.DatabaseHandler.UpdateConnectionStatus(1, 1)
		if err != nil {
			fmt.Println("Error updating connection status:", err)
		}
		fmt.Println("Client connected.")
	case *events.Disconnected:
		err := meowController.DatabaseHandler.UpdateConnectionStatus(1, 0)
		if err != nil {
			fmt.Println("Error updating connection status:", err)
		}
		fmt.Println("Client disconnected.")
	}
}

func (meowController *WhatsMeowController) Connect() {
	go meowController.StartClient()
}

func (meowController *WhatsMeowController) Disconnect() {
	meowController.Client.Disconnect()
	err := meowController.DatabaseHandler.UpdateConnectionStatus(1, 0)
	if err != nil {
		fmt.Println("Error updating connection status:", err)
	}
}

func (meowController *WhatsMeowController) Logout() {
	if meowController.Client == nil {
		fmt.Println("Client pointer is nil.")
		return
	} else {
		if meowController.Client.IsLoggedIn() && meowController.Client.IsConnected() {
			meowController.Client.Logout()
		} else {
			fmt.Println("Client not logged in or connected.")
		}
	}
}

func (meowController *WhatsMeowController) GetStatus() (isConnected, isLoggedIn bool, err error) {
	if meowController.Client == nil {
		return false, false, fmt.Errorf("client pointer is nil")
	}
	isConnected = meowController.Client.IsConnected()
	isLoggedIn = meowController.Client.IsLoggedIn()
	return isConnected, isLoggedIn, nil
}

func (meowController *WhatsMeowController) GetQr() (qr string) {
	if meowController.Client == nil {
		return
	}
	if !meowController.Client.IsConnected() {
		return
	}
	if meowController.Client.IsLoggedIn() {
		return
	}
	qrCode, err := meowController.DatabaseHandler.GetQrCode()
	if err != nil {
		return
	}
	return qrCode
}

func (meowController *WhatsMeowController) SendTextMessage(phone string, message string) {
	if phone[0:1] == "+" { // Remove the "+" sign if it exists
		phone = phone[1:]
	}
	if phone[0:1] == "0" { // Replace the first "0" with "62" if it starts with "0"
		phone = "62" + phone[1:]
	}
	// remove any non-numeric characters
	phone = regexp.MustCompile("[^0-9]+").ReplaceAllString(phone, "")

	// send message
	meowController.Client.SendMessage(context.Background(), types.JID{User: phone, Server: "s.whatsapp.net"}, &waProto.Message{Conversation: proto.String(message)})
}

func (meowController *WhatsMeowController) SendTextMessageGroup(jid string, message string) {
	meowController.Client.SendMessage(context.Background(), types.JID{User: jid, Server: "g.us"}, &waProto.Message{
		Conversation: proto.String(message),
	})
}
