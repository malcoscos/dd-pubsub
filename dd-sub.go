package dd_pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	websocket "github.com/gorilla/websocket"
	minio "github.com/minio/minio-go/v7"
	credentials "github.com/minio/minio-go/v7/pkg/credentials"
)

var ctx = context.Background()

func Subscribe(s *SubArg) {

	// make channel
	msgCh := make(chan mqtt.Message)
	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		msgCh <- msg
	}
	// configure mqtt client options
	opts := mqtt.NewClientOptions()
	//　add broker
	broker := fmt.Sprintf("tcp://%s:%s", "127.0.0.1", "1883")
	opts.AddBroker(broker)
	// make mqtt client
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
	// notify systemcall
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

			if descriptor.DataType == "video_data" {
				url := fmt.Sprintf("ws://%s/%s", descriptor.DatabaseAddr, descriptor.Locator)
				dialer := websocket.DefaultDialer
				// WebSocketサーバーに接続
				conn, _, err := dialer.Dial(url, nil)

				if err != nil {
					log.Fatalf("WebSocket接続に失敗: %v", err)
				}
				defer conn.Close()

				// サーバーからのメッセージを受信して表示するループ
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Fatalf("メッセージの読み取りに失敗: %v", err)
				}
				// 受信したメッセージを表示
				fmt.Printf("受信メッセージ: %s\n", message)

			} else if descriptor.DataType == "image" || descriptor.DataType == "tiny_data" {
				database_addr := fmt.Sprintf("%s:%s", descriptor.DatabaseAddr, descriptor.DatabasePort)
				accessKeyID := "hoge"          // アクセスキーID
				secretAccessKey := "hoge_hoge" // シークレットアクセスキー
				useSSL := false                // SSLを使用する場合はtrueに設定

				// MinIOクライアントの初期化
				minioClient, err := minio.New(database_addr, &minio.Options{
					Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
					Secure: useSSL,
				})
				if err != nil {
					log.Fatalln(err)
				}

				// オブジェクトを取得
				bucket_name := descriptor.Topic
				object_name := descriptor.Locator
				object, err := minioClient.GetObject(context.Background(), bucket_name, object_name, minio.GetObjectOptions{})
				if err != nil {
					fmt.Print("helllo")
					log.Fatalln(err)
				}
				log.Printf("Successfully download %s", object_name)
				defer object.Close()
			}

		// to interrupt if there is systemcall
		case <-signalCh:
			fmt.Printf("Interrupt detected.\n")
			c.Disconnect(1000)
			return
		}
	}
}
