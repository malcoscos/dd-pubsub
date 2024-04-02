package retrieve_func

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/malcoscos/dd-pubsub/types"
	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
)

func RetreiveTinyData(descriptor types.Descriptor, s *types.SubArg) {
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

	// オブジェクトのデータを読み込む
	data, err := io.ReadAll(object)
	if err != nil {
		log.Fatalln(err)
	}

	// create file
	file_path := fmt.Sprintf("%s/%s", s.StorePath, descriptor.TimeStamp)
	file, err := os.Create(file_path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// write file
	_, err = file.Write(data)
	if err != nil {
		fmt.Println(err)
	}
}
