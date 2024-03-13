package dd_pubsub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	uuid "github.com/google/uuid"
	minio "github.com/minio/minio-go/v7"
	credentials "github.com/minio/minio-go/v7/pkg/credentials"
)

func Publish(p *PubArg) {
	var data_mime_type string = ProcessFile(p.Payload)

	object_name := uuid.NewString()

	if data_mime_type == "video" {

		// データをファイルシステムに保存
		dir := "./saved_files"
		payload_data, ok := p.Payload.([]byte)
		if !ok {
			return
		}
		// データをファイルに保存し、ファイルパスを取得
		filepath, err := saveDataToFile(payload_data, dir, object_name)
		if err != nil {
			fmt.Println("Failed to save data to file:", err)
			return
		}
		fmt.Println("Data saved to file:", filepath)
		object_name = filepath

	} else if data_mime_type == "image" || data_mime_type == "tiny_data" {
		// configure minio addr and auth
		database_addr := fmt.Sprintf("%s:%s", p.StrageAddr, p.StragePort)
		useSSL := false //recommend to change to true in the production env
		// create minio client
		var minioClient *minio.Client
		var err error
		if p.StrageId != "" {
			// authenticated communication
			minioClient, err = minio.New(database_addr, &minio.Options{
				Creds:  credentials.NewStaticV4(p.StrageId, p.StrageKey, ""),
				Secure: useSSL,
			})
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			// anonymous communication
			minioClient, err = minio.New(database_addr, &minio.Options{
				Creds:  credentials.NewIAM(""),
				Secure: useSSL,
			})
			if err != nil {
				log.Fatalln(err)
			}
		}
		// upload data to minio
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
	}

	// descriptor of real data
	now := time.Now()
	time_stamp := fmt.Sprint(now.Format(time.RFC3339))
	descriptor := Descriptor{
		Topic:        p.Topic,
		DatabaseAddr: p.StrageAddr,
		DatabasePort: p.StragePort,
		DataType:     data_mime_type,
		Locator:      object_name,
		TimeStamp:    time_stamp,
		Header:       "hoge", // This attr is used after the ffmpeg implementation is finished
	}

	// to encode from golang structure to json
	jsonData, err := json.Marshal(descriptor)
	if err != nil {
		fmt.Println(err)
		return
	}

	// publich to broker
	token := p.MqttClient.Publish(p.Topic, p.Qos, p.Retained, jsonData)
	token.Wait()
	fmt.Println("Complete publish")
}

func store_tiny_data() {

}

func store_video_data() {

}

func saveDataToFile(data []byte, dir, file_name string) (string, error) {
	// ディレクトリを作成（存在しない場合）
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// ファイル名を生成（プレフィックス+タイムスタンプ）
	fullPath := filepath.Join(dir, file_name)

	// ファイルを開く（存在しない場合は作成、存在する場合は上書き）
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// データをファイルに書き込む
	if _, err = file.Write(data); err != nil {
		return "", err
	}

	return fullPath, nil
}
