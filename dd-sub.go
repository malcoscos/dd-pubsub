package main

import (
	"fmt"
	"ftp"
	"log"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// channelの作成
	msgCh := make(chan mqtt.Message)
	// messageをchannelに送信する関数の作成
	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		msgCh <- msg
	}
	// optsにClientOptionsインスタンスのpointerを格納
	opts := mqtt.NewClientOptions()
	//　BrokerServerのlistに追加
	opts.AddBroker("tcp://localhost:1883")
	// clientクラスのインスタンスを作成
	c := mqtt.NewClient(opts)
	// Brokerへのconnection及び、Errorがないか判定
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Mqtt error: %s", token.Error())
	}
	// subscribeするtopicが、正しく設定されているか判定
	if subscribeToken := c.Subscribe("go-mqtt/sample", 0, f); subscribeToken.Wait() && subscribeToken.Error() != nil {
		log.Fatal(subscribeToken.Error())
	}

	// systemcallを受け取るchanenlの作成
	signalCh := make(chan os.Signal, 1)
	// systemcallがあると知らせる
	signal.Notify(signalCh, os.Interrupt)
	// forever
	for {
		select {
		// メッセージをchanellから受信
		case m := <-msgCh:
			// fmt.Printf("topic: %v, payload: %v\n", m.Topic(), string(m.Payload()))
			file_name := string(m.Payload())
			data, err := os.ReadFile(string(file_name))
			if err != nil {
				log.Fatal(err)
			}
			println(string(data))

			ftp.Connect("10.0.8.19", 21)
			download_file := ftp.Download(file_name, dest)

		// systemcallがあると知らせる
		case <-signalCh:
			fmt.Printf("Interrupt detected.\n")
			c.Disconnect(1000)
			return
		}
	}
}
