package dd_pubsub

// if descriptor.DataType == "tiny" {
// 	// info of data
// 	redis_addr := fmt.Sprintf("%s:%s", descriptor.DatabaseAddr, descriptor.DatabasePort)
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     redis_addr, // Redisサーバーのアドレス
// 		Password: "",         // パスワードがない場合は空文字列
// 		DB:       0,          // 使用するデータベース
// 	})

// 	// キーから値を取得
// 	val, err := rdb.Get(ctx, descriptor.Locator).Result()
// 	if err != nil {
// 		log.Fatalf("Failed to get key: %v", err)
// 	}
// 	fmt.Println("get this data: ", val)

// 	// Redisクライアントのクローズ
// 	err = rdb.Close()
// 	if err != nil {
// 		log.Fatalf("Failed to close client: %v", err)
// 	}
// } else {
// 	database_addr := fmt.Sprintf("%s:%s", descriptor.DatabaseAddr, descriptor.DatabasePort)
// 	accessKeyID := "hoge"          // アクセスキーID
// 	secretAccessKey := "hoge_hoge" // シークレットアクセスキー
// 	useSSL := false                // SSLを使用する場合はtrueに設定

// 	// MinIOクライアントの初期化
// 	minioClient, err := minio.New(database_addr, &minio.Options{
// 		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
// 		Secure: useSSL,
// 	})
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	bucket_name := descriptor.Topic   // バケット名
// 	object_name := descriptor.Locator // オブジェクト名

// 	// オブジェクトを取得
// 	object, err := minioClient.GetObject(context.Background(), bucket_name, object_name, minio.GetObjectOptions{})
// 	if err != nil {
// 		fmt.Print("helllo")
// 		log.Fatalln(err)
// 	}
// 	log.Printf("Successfully download %s", object_name)
// 	defer object.Close()
// }
