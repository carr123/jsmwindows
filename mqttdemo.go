package main

import (
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	go NatsMQTTRecv()
	time.Sleep(time.Second)
	go NatsMQTTSend()
	time.Sleep(time.Hour * 1000)
}

func NatsMQTTRecv() {
	BrokerURL := "127.0.0.1"
	ClientID := "mqtt_recv"
	username := "iot"
	password := "abc123"
	SubTopic := "Client/Device/202501"

	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://" + BrokerURL + ":1883")
	opts.SetClientID(ClientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(time.Second * 3)
	opts.SetResumeSubs(true)
	opts.SetCleanSession(false)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		log.Println("recv:", msg.Topic(), " ", msg.Payload())
	})

	MQClient := mqtt.NewClient(opts)
	token := MQClient.Connect()
	token.Wait()
	if token.Error() != nil {
		log.Println("recv Connect err:", token.Error())
		return
	}

	tk := MQClient.Subscribe(SubTopic, 1, func(client mqtt.Client, msg mqtt.Message) {
		data := string(msg.Payload())
		log.Println("sub recv:", " msgid:", msg.MessageID(),
			" topic:", msg.Topic(), " data:", data, " retain:", msg.Retained())
	})
	tk.Wait()
	if tk.Error() != nil {
		log.Println("Subscribe err:", token.Error())
		return
	}

	log.Println("begin sub msg...")

	time.Sleep(time.Hour * 8000)
}

func NatsMQTTSend() {
	BrokerURL := "127.0.0.1"
	ClientID := "mqtt_send"
	username := "iot"
	password := "abc123"

	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://" + BrokerURL + ":1883")
	opts.SetClientID(ClientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(time.Second * 3)
	opts.SetResumeSubs(true)
	opts.SetCleanSession(false)

	MQClient := mqtt.NewClient(opts)
	token := MQClient.Connect()
	token.Wait()
	if token.Error() != nil {
		log.Println("sender Connect:", token.Error())
		return
	}

	for {
		time.Sleep(time.Second * 2)
		content := time.Now().String()

		tk := MQClient.Publish("Client/Device/202501", 1, false, content)
		tk.Wait()
		if tk.Error() != nil {
			log.Println("Publish err:", tk.Error())
			continue
		}
		log.Println("send:", content)
	}

	log.Println("after send")

	time.Sleep(time.Hour * 8000)
}
