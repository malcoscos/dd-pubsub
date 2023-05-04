package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	auth "github.com/bramvdbogaerde/go-scp/auth"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	scp "github.com/lkbhargav/go-scp"
	ssh "golang.org/x/crypto/ssh"
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
	opts.AddBroker("tcp://10.0.255.76:1883")
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

			// Use SSH key authentication from the auth package
			// we ignore the host key in this example, please change this if you use this library
			clientConfig, _ := auth.PrivateKey("shinoda-lab", "~/.ssh", ssh.InsecureIgnoreHostKey())

			// Create a new SCP client
			client := scp.NewClient("10.0.8.19:21", &clientConfig)

			// Connect to the remote server
			err_connect := client.Connect()

			if err_connect != nil {
				fmt.Println("Couldn't establish a connection to the remote server ", err_connect)
				return
			}

			// Open a file
			f, _ := os.Open("/tmp")

			// Close client connection after the file has been copied
			defer client.Close()

			// Close the file after it has been copied
			defer f.Close()

			// Finaly, copy the file over
			// Usage: CopyFile(fileReader, remotePath, permission)

			err_copy_file := client.CopyFile(f, "file_name", "0655")

			if err_copy_file != nil {
				fmt.Println("Error while copying file ", err_copy_file)
			}

		// systemcallがあると知らせる
		case <-signalCh:
			fmt.Printf("Interrupt detected.\n")
			c.Disconnect(1000)
			return
		}
	}
}
