package dd_pubsub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// json format
type Descriptor struct {
	// Addr    string
	// Port    string
	Format  string
	Locator string
}

type PubArg struct {
	Topic         string
	Qos           byte
	Retained      bool
	Payload       interface{}
	DataFormat    string
	BrokerAddr    string
	BrokerPort    string
	NFSServerAddr string
	NFSServerPort string
	FilePath      string
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

	// JSON形式のデータを[]byteに変換
	payloadBytes, err := json.Marshal(p.Payload)
	if err != nil {
		log.Fatalf("Error marshalling payload: %s", err)
	}

	// ファイルにデータを書き込む
	err = os.WriteFile(p.FilePath, payloadBytes, 0644)
	if err != nil {
		log.Fatalf("Error writing file: %s", err)
	}

	// info of data
	payload_data := Descriptor{
		// Addr:    p.NFSServerAddr,
		// Port:    p.NFSServerPort,
		Format:  p.DataFormat,
		Locator: p.FilePath,
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

	c.Disconnect(250)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Complete publish")
}
