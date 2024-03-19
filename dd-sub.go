package dd_pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	websocket "github.com/gorilla/websocket"
	types "github.com/malcoscos/dd-pubsub/types"
	minio "github.com/minio/minio-go/v7"
	credentials "github.com/minio/minio-go/v7/pkg/credentials"
)

func Subscribe(s *types.SubArg) {

	// make channel
	msgCh := make(chan mqtt.Message)
	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		msgCh <- msg
	}
	// mqtt client
	c := s.MqttClient
	//connect to broker
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	// subscribe from broker
	if subscribeToken := c.Subscribe(s.Topic, s.Qos, f); subscribeToken.Wait() && subscribeToken.Error() != nil {
		fmt.Println(subscribeToken.Error())
	}

	//  notify systemcall get from channel
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	// forever
	for {
		select {
		// get message from channel
		case m := <-msgCh:
			var descriptor types.Descriptor
			payload_data := string(m.Payload())

			// to decode from golong structure to json
			if err := json.Unmarshal([]byte(payload_data), &descriptor); err != nil {
				fmt.Println(err)
				return
			}

			if descriptor.DataType == "video_data" {
				// connect to server of websocket
				url := fmt.Sprintf("ws://%s/%s", descriptor.DatabaseAddr, descriptor.Locator)
				dialer := websocket.DefaultDialer
				conn, _, err := dialer.Dial(url, nil)
				if err != nil {
					fmt.Println(err)
				}
				defer conn.Close()

				// サーバーからのメッセージを受信して表示するループ
				_, message, err := conn.ReadMessage()
				if err != nil {
					fmt.Println(err)
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
					fmt.Println(err)
				}

				// オブジェクトを取得
				bucket_name := descriptor.Topic
				object_name := descriptor.Locator
				object, err := minioClient.GetObject(context.Background(), bucket_name, object_name, minio.GetObjectOptions{})
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Successfully download %s", object_name)
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
