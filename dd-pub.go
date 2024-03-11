package dd_pubsub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	minio "github.com/minio/minio-go/v7"
	credentials "github.com/minio/minio-go/v7/pkg/credentials"
)

func Publish(p *PubArg) {
	var data_mime_type string = ProcessFile(p.Payload)

	// configure minio addr and auth
	database_addr := fmt.Sprintf("%s:%s", p.StrageAddr, p.StragePort)
	accessKeyID := "hoge"
	secretAccessKey := "hoge_hoge"
	useSSL := false

	// create minio client
	var minioClient *minio.Client
	var err error

	if p.StrageId != "" {
		minioClient, err = minio.New(database_addr, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		minioClient, err = minio.New(database_addr, &minio.Options{
			Creds:  credentials.Anonymous(),
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
	object_name := uuid.NewString()
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

	// get timestamp
	now := time.Now()
	time_stamp := fmt.Sprint(now.Format(time.RFC3339))

	// info of data
	descriptor := Descriptor{
		Topic:        p.Topic,
		DatabaseAddr: p.StrageAddr,
		DatabasePort: p.StragePort,
		DataType:     data_mime_type,
		Locator:      object_name,
		TimeStamp:    time_stamp,
		Header:       "hoge",
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

	// mqttクライアントのクローズ
	fmt.Println("Complete publish")
}
