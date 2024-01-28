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
)

// json format
type Descriptor struct {
	Format       string
	Locator      string
	DatabaseAddr string
	DatabasePort string
	TimeStamp    string
	Header       string
}

type PubArg struct {
	Topic        string
	Qos          byte
	Retained     bool
	Payload      interface{}
	DataFormat   string
	BrokerAddr   string
	BrokerPort   string
	DatabaseAddr string
	DatabasePort string
}

func Publish(p *PubArg) {

	// ClientOptionsインスタンスのpointerを格納
	opts := mqtt.NewClientOptions()

	//　add broker
	broker := fmt.Sprintf("tcp://%s:%s", p.BrokerAddr, p.BrokerPort)
	opts.AddBroker(broker)

	// clientのインスタンスを作成
	c := mqtt.NewClient(opts)

	// connect to broker
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Mqtt error: %s", token.Error())
	}

	var ctx = context.Background()

	// Redisクライアントの作成
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

	// Redisクライアントのクローズ
	err = rdb.Close()
	if err != nil {
		log.Fatalf("Failed to close client: %v", err)
	}
}
