package dd-pubsub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// json format
type Payload struct {
	Addr     string
	Port     string
	Format   string
	Location string
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
	broker := fmt.Sprintf("tcp://%d:%d", p.BrokerAddr, p.BrokerPort)
	opts.AddBroker(broker)

	// clientのインスタンスを作成
	c := mqtt.NewClient(opts)

	// connect to broker
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Mqtt error: %s", token.Error())
	}

	// to write file in pub server
	err := os.WriteFile(p.FilePath, p.Payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	// info of data
	payload_data := Payload{
		Addr:     p.NFSServerAddr,
		Port:     p.NFSServerPort,
		Format:   p.DataFormat,
		Location: p.FilePath
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
