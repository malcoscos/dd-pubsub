package dd_pubsub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	minio "github.com/minio/minio-go/v7"
	credentials "github.com/minio/minio-go/v7/pkg/credentials"
)

func Publish(p *PubArg) {

	// ClientOptionsインスタンスのpointerを格納
	opts := mqtt.NewClientOptions()

	//　add broker
	broker := fmt.Sprintf("tcp://%s:%s", p.BrokerAddr, p.BrokerPort)
	opts.AddBroker(broker)

	// clientのインスタンスを作成
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Mqtt error: %s", token.Error())
	}

	// データに応じて送信先の分岐
	var data_mime_type string = ProcessFile(p.Payload)

	var key string

	if data_mime_type == "tiny" {
		// redisに対して送信

		// Redisクライアントの作成
		var ctx = context.Background()
		database_addr := fmt.Sprintf("%s:%s", p.RedisAddr, p.RedisPort)
		rdb := redis.NewClient(&redis.Options{
			Addr:     database_addr, // Redisサーバーのアドレス
			Password: "",            // パスワードがない場合は空文字列
			DB:       0,             // 使用するデータベース
		})

		// 生成されたIDをキーとして使用
		key := fmt.Sprintf("%s:%s", p.Topic, uuid.NewString())

		// キーと値をセット
		err := rdb.Set(ctx, key, p.Payload, 0).Err()
		if err != nil {
			log.Fatalf("Error setting value: %v", err)
		}

		// Redisクライアントのクローズ
		err = rdb.Close()
		if err != nil {
			log.Fatalf("Failed to close client: %v", err)
		}
	} else {
		// オブジェクトストレージに対して送信
		database_addr := fmt.Sprintf("%s:%s", p.MinioAddr, p.MinioPort) // MinIOサーバーのアドレスとポート
		accessKeyID := "hoge"                                           // アクセスキー
		secretAccessKey := "hoge_hoge"                                  // シークレットキー
		useSSL := false                                                 // SSLを使用する場合はtrueに設定

		// MinIOクライアントの初期化
		minioClient, err := minio.New(database_addr, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			log.Fatalln(err)
		}

		// オブジェクトストレージにアップロード
		bucket_name := p.Topic
		exists, err := minioClient.BucketExists(ctx, bucket_name)
		if err != nil {
			log.Fatalln(err)
		}
		if !exists {
			err = minioClient.MakeBucket(context.Background(), bucket_name, minio.MakeBucketOptions{})
			if err != nil {
				log.Fatalln(err)
			}
		}
		uuid := uuid.NewString()
		object_name := uuid
		payload_data, ok := p.Payload.([]byte)
		if !ok {
			return
		}
		var reader io.Reader
		if data, ok := p.Payload.([]byte); ok {
			reader = bytes.NewReader(data)
		} else {
			log.Fatalln("Payload is not of type []byte")
		}
		info, err := minioClient.PutObject(context.Background(), bucket_name, object_name, reader, int64(len(payload_data)), minio.PutObjectOptions{})
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Successfully uploaded %s of size %d\n", object_name, info.Size)
		key = object_name
	}

	// 時刻の取得
	now := time.Now()
	time_stamp := fmt.Sprint(now.Format(time.RFC3339))

	// info of data
	payload_data := Descriptor{
		Topic:        p.Topic,
		DatabaseAddr: p.RedisAddr,
		DatabasePort: p.RedisPort,
		DataType:     data_mime_type,
		Locator:      key,
		TimeStamp:    time_stamp,
		Header:       "hoge",
	}

	// to encode from golang structure to json
	jsonData, err := json.Marshal(payload_data)
	if err != nil {
		fmt.Println(err)
		return
	}

	// publich to broker
	token := c.Publish(p.Topic, p.Qos, p.Retained, jsonData)
	token.Wait()

	// mqttクライアントのクローズ
	c.Disconnect(250)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Complete publish")
}

func ProcessFile(data interface{}) string {
	bytes, ok := data.([]byte)
	if !ok {
		return "huge"
	}

	mimeType := http.DetectContentType(bytes)

	switch {
	case mimeType == "video/mp4" || mimeType == "video/x-msvideo":
		fmt.Println("Processing video file with ffmpeg...")
		// ここでffmpegを実行します。
		// 実際にはffmpegコマンドはファイルを要求するため、ファイルへの書き出しが必要です。
		return "huge"
	case mimeType == "image/jpeg" || mimeType == "image/png":
		fmt.Println("Processing image file with jhead...")
		// ここでjheadを実行します。
		// 実際にはjheadコマンドもファイルを要求するため、ファイルへの書き出しが必要です。
		return "huge"
	default:
		fmt.Println("No action required for this file type.")
		return "tiny"
	}
}
