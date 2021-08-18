package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	nats "github.com/nats-io/nats.go"
)

var (
	servers = []string{"nats://127.0.0.1:4222", "nats://127.0.0.1:5222", "nats://127.0.0.1:6222"}
	nats_user   = "your_user"
	nats_passwd = "your_passwd"
)

func main() {
	jsmSend()
}

func jsmSend() {
	var err error

	nc, err := nats.Connect(
		strings.Join(servers, ","),
		nats.Name("jsmSend"),
		nats.Timeout(10*time.Second),
		nats.PingInterval(20*time.Second),
		nats.MaxPingsOutstanding(5),
		nats.MaxReconnects(10),
		nats.ReconnectWait(10*time.Second),
		nats.ReconnectBufSize(5*1024*1024),
		nats.UserInfo(nats_user, nats_passwd),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "ORDERS",
		Subjects:  []string{"ORDERS.*"},
		Storage:   nats.FileStorage,
		Replicas:  3,
		Retention: nats.WorkQueuePolicy,
		Discard:   nats.DiscardNew,
		MaxMsgs:   -1,
		MaxAge:    time.Hour * 24 * 365,
	})
	if err != nil && !strings.Contains(err.Error(), "already in use") {
		log.Println("2", err)
		return
	}

	for i := 0; i < 5000; i++ {
		szMsg := fmt.Sprintf("js msg %02d", i)
		js.PublishAsync("ORDERS.created", []byte(szMsg))
	}
	<-js.PublishAsyncComplete()

	log.Println("sender complete. count =", 5000)
	time.Sleep(time.Second * 10)
	log.Println("sender exit")
}
