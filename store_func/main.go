package store_func

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/malcoscos/dd-pubsub/types"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func StoreTinyData(p *types.PubArg, object_name string) string {
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
	var ctx = context.Background()
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
	var object_size int64

	if tiny_payload_data, ok := p.Payload.([]byte); ok {
		reader = bytes.NewReader(tiny_payload_data)
		object_size = int64(len(tiny_payload_data))
	} else {
		log.Fatalln("Payload is not of type []byte")
	}
	info, err := minioClient.PutObject(context.Background(), bucket_name, object_name, reader, object_size, minio.PutObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", object_name, info.Size)
	return object_name
}

func StoreVideoData(data interface{}, object_name string, dir string) string {
	payload_data, ok := data.([]byte)
	if !ok {
		fmt.Println("Failed to exchange data to byte", ok)
		return ""
	}
	// データをファイルに保存し、ファイルパスを取得
	filepath, err := SaveDataToFile(payload_data, dir, object_name)
	if err != nil {
		fmt.Println("Failed to save data to file:", err)
		return ""
	}
	fmt.Println("Data saved to file:", filepath)
	return filepath
}

func SaveDataToFile(data []byte, dir, file_name string) (string, error) {
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
