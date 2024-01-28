package dd_pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	redis "github.com/go-redis/redis/v8"
)

type SubArg struct {
	Topic      string
	Qos        byte
	BrokerAddr string
	BrokerPort string
}

var ctx = context.Background()

func Subscribe(s *SubArg) {

	// channelの作成
	msgCh := make(chan mqtt.Message)

	// messageをchannelに送信する関数の作成
	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		msgCh <- msg
	}

	// optsにClientOptionsインスタンスのpointerを格納
	opts := mqtt.NewClientOptions()

	//　add broker
	broker := fmt.Sprintf("tcp://%s:%s", s.BrokerAddr, s.BrokerPort)
	opts.AddBroker(broker)

	// clientのインスタンスを作成
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
			redis_addr := fmt.Sprintf("%s:%s", descriptor.DatabaseAddr, descriptor.DatabasePort)
			rdb := redis.NewClient(&redis.Options{
				Addr:     redis_addr, // Redisサーバーのアドレス
				Password: "",         // パスワードがない場合は空文字列
				DB:       0,          // 使用するデータベース
			})

			// キーから値を取得
			val, err := rdb.Get(ctx, descriptor.Locator).Result()
			if err != nil {
				log.Fatalf("Failed to get key: %v", err)
			}
			fmt.Println("get this data: ", val)

			// Redisクライアントのクローズ
			err = rdb.Close()
			if err != nil {
				log.Fatalf("Failed to close client: %v", err)
			}

		// to interrupt if there is systemcall
		case <-signalCh:
			fmt.Printf("Interrupt detected.\n")
			c.Disconnect(1000)
			return
		}
	}
}
