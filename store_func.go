package dd_pubsub

// func branch_dst_strage() {
// 	// Branching of destination depending on data
// 	var data_mime_type string = ProcessFile(p.Payload)

// 	if data_mime_type == "" {
// 		// redisに対して送信

// 		// Redisクライアントの作成
// 		var ctx = context.Background()
// 		database_addr := fmt.Sprintf("%s:%s", p.RedisAddr, p.RedisPort)
// 		rdb := redis.NewClient(&redis.Options{
// 			Addr:     database_addr, // Redisサーバーのアドレス
// 			Password: "",            // パスワードがない場合は空文字列
// 			DB:       0,             // 使用するデータベース
// 		})

// 		// 生成されたIDをキーとして使用
// 		key := fmt.Sprintf("%s:%s", p.Topic, uuid.NewString())

// 		// キーと値をセット
// 		err := rdb.Set(ctx, key, p.Payload, 0).Err()
// 		if err != nil {
// 			log.Fatalf("Error setting value: %v", err)
// 		}

// 		// Redisクライアントのクローズ
// 		err = rdb.Close()
// 		if err != nil {
// 			log.Fatalf("Failed to close client: %v", err)
// 		}
// 	} else {
// 		// オブジェクトストレージに対して送信
// 		database_addr := fmt.Sprintf("%s:%s", p.MinioAddr, p.MinioPort) // MinIOサーバーのアドレスとポート
// 		accessKeyID := "hoge"                                           // アクセスキー
// 		secretAccessKey := "hoge_hoge"                                  // シークレットキー
// 		useSSL := false                                                 // SSLを使用する場合はtrueに設定

// 		// MinIOクライアントの初期化
// 		minioClient, err := minio.New(database_addr, &minio.Options{
// 			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
// 			Secure: useSSL,
// 		})
// 		if err != nil {
// 			log.Fatalln(err)
// 		}

// 		// オブジェクトストレージにアップロード
// 		bucket_name := p.Topic
// 		exists, err := minioClient.BucketExists(ctx, bucket_name)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 		if !exists {
// 			err = minioClient.MakeBucket(context.Background(), bucket_name, minio.MakeBucketOptions{})
// 			if err != nil {
// 				log.Fatalln(err)
// 			}
// 		}
// 		uuid := uuid.NewString()
// 		object_name := uuid
// 		payload_data, ok := p.Payload.([]byte)
// 		if !ok {
// 			return
// 		}
// 		var reader io.Reader
// 		if data, ok := p.Payload.([]byte); ok {
// 			reader = bytes.NewReader(data)
// 		} else {
// 			log.Fatalln("Payload is not of type []byte")
// 		}
// 		info, err := minioClient.PutObject(context.Background(), bucket_name, object_name, reader, int64(len(payload_data)), minio.PutObjectOptions{})
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 		log.Printf("Successfully uploaded %s of size %d\n", object_name, info.Size)
// 		key = object_name
// 	}
// }
