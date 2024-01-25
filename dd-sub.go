package dd_pubsub

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	auth "github.com/bramvdbogaerde/go-scp/auth"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	scp "github.com/povsister/scp"
	ssh "golang.org/x/crypto/ssh"
)

type SubArg struct {
	Topic           string
	Qos             byte
	BrokerAddr      string
	BrokerPort      string
	NFSServerAddr   string
	NFSServerPort   string
	SSHUsername     string
	SSHPassword     string
	CopyFileDstPath string
}

func Subscribe(s *SubArg) {
	// channelの作成
	msgCh := make(chan mqtt.Message)
	// messageをchannelに送信する関数の作成
	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		msgCh <- msg
	}
	// optsにClientOptionsインスタンスのpointerを格納
	opts := mqtt.NewClientOptions()

	//　add broker to list
	broker := fmt.Sprintf("tcp://%s:%s", s.BrokerAddr, s.BrokerPort)
	opts.AddBroker(broker)

	// make client instance
	c := mqtt.NewClient(opts)

	//connect to broker
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Mqtt error: %s", token.Error())
	}

	// subscribe from broker
	if subscribeToken := c.Subscribe(s.Topic, s.Qos, f); subscribeToken.Wait() && subscribeToken.Error() != nil {
		log.Fatal(subscribeToken.Error())
	}

	// systemcallを受け取るchanenlの作成
	signalCh := make(chan os.Signal, 1)
	// systemcallがあると知らせる
	signal.Notify(signalCh, os.Interrupt)

	// forever
	for {
		select {
		// get message from channel
		case m := <-msgCh:
			var descriptor Descriptor
			payload_data := string(m.Payload())

			// to decode from golong structure to json
			if err := json.Unmarshal([]byte(payload_data), &descriptor); err != nil {
				fmt.Println(err)
				return
			}

			// info of data
			file_name_nfs := descriptor.Locator
			server_addr := fmt.Sprintf("%s:%s", s.NFSServerAddr, s.NFSServerPort)

			// auth and create a new SCP client
			client_config, _ := auth.PasswordKey(s.SSHUsername, s.SSHPassword, ssh.InsecureIgnoreHostKey())
			client, err_connect := scp.NewClient(server_addr, &client_config, &scp.ClientOption{})

			// Connect to the remote server
			if err_connect != nil {
				fmt.Println("Couldn't establish a connection to the remote server ", err_connect)
				return
			}

			// copy the file over
			err_copy_file := client.CopyFileFromRemote(file_name_nfs, s.CopyFileDstPath, &scp.FileTransferOption{})
			if err_copy_file != nil {
				fmt.Println("Error while copying file ", err_copy_file)
			}

		// to interrupt if there is systemcall
		case <-signalCh:
			fmt.Printf("Interrupt detected.\n")
			c.Disconnect(1000)
			return
		}
	}
}
