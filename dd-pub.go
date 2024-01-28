package dd_pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	redis "github.com/go-redis/redis/v8"
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
	var data_mime_type bool = ProcessFile(p.Payload)

	var key string

	if data_mime_type {
		// redisに対して送信

		// Redisクライアントの作成
		var ctx = context.Background()
		redis_addr := fmt.Sprintf("%s:%s", p.DatabaseAddr, p.DatabasePort)
		rdb := redis.NewClient(&redis.Options{
			Addr:     redis_addr, // Redisサーバーのアドレス
			Password: "",         // パスワードがない場合は空文字列
			DB:       0,          // 使用するデータベース
		})

		// idを作成
		newId, err := rdb.Incr(ctx, "unique_counter").Result()
		if err != nil {
			log.Fatalf("Error incrementing counter: %v", err)
		}

		// 生成されたIDをキーとして使用
		key := fmt.Sprintf("%s:%d", p.Topic, newId)

		// キーと値をセット
		err = rdb.Set(ctx, key, p.Payload, 0).Err()
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
		endpoint := "your-minio-endpoint:9000" // MinIOサーバーのアドレスとポート
		accessKeyID := "your-access-key"       // アクセスキー
		secretAccessKey := "your-secret-key"   // シークレットキー
		useSSL := false                        // SSLを使用する場合はtrueに設定

		// MinIOクライアントの初期化
		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			log.Fatalln(err)
		}

		bucketName := "your-bucket-name" // バケット名
		objectName := "your-object-name" // アップロードするオブジェクトの名前
		filePath := "your-file-path"     // アップロードするファイルのパス

		// ファイルをアップロード
		info, err := minioClient.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: "application/octet-stream"})
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	}

	// 時刻の取得
	now := time.Now()
	time_stamp := fmt.Sprint(now.Format(time.RFC3339))

	// info of data
	payload_data := Descriptor{
		DatabaseAddr: p.DatabaseAddr,
		DatabasePort: p.DatabasePort,
		Format:       p.DataFormat,
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

func ProcessFile(data interface{}) bool {
	bytes, ok := data.([]byte)
	if !ok {
		return true
	}

	mimeType := http.DetectContentType(bytes)

	switch {
	case mimeType == "video/mp4" || mimeType == "video/x-msvideo":
		fmt.Println("Processing video file with ffmpeg...")
		// ここでffmpegを実行します。
		// 実際にはffmpegコマンドはファイルを要求するため、ファイルへの書き出しが必要です。
		return false
	case mimeType == "image/jpeg" || mimeType == "image/png":
		fmt.Println("Processing image file with jhead...")
		// ここでjheadを実行します。
		// 実際にはjheadコマンドもファイルを要求するため、ファイルへの書き出しが必要です。
		return false
	default:
		fmt.Println("No action required for this file type.")
		return true
	}
}
