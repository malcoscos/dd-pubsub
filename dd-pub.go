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

func Publish(topic string, qos byte, retained bool, payload interface{}, data_format string, broker_addr string, broker_port string, nfs_server_addr string, nfs_server_port string, file_path string) {

	// ClientOptionsインスタンスのpointerを格納
	opts := mqtt.NewClientOptions()

	//　add broker
	broker := fmt.Sprintf("tcp://%d:%d", broker_addr, broker_port)
	opts.AddBroker(broker)

	// clientのインスタンスを作成
	c := mqtt.NewClient(opts)

	// connect to broker
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Mqtt error: %s", token.Error())
	}

	// to write file in pub server
	err := os.WriteFile(file_path, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	// info of data
	payload_data := Payload{
		Addr:     nfs_server_addr,
		Port:     nfs_server_port,
		Format:   data_format,
		Location: file_path
	}

	// to encode from golang structure to json
	jsonData, err := json.Marshal(payload_data)
	if err != nil {
		fmt.Println(err)
		return
	}

	// publich to broker
	token := c.Publish(topic, pos, retained, jsonData)
	token.Wait()

	c.Disconnect(250)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Complete publish")
}
