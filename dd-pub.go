package main

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

func main() {

	// ClientOptionsインスタンスのpointerを格納
	opts := mqtt.NewClientOptions()

	//　add broker
	opts.AddBroker("tcp://10.0.8.25:1883")

	// clientのインスタンスを作成
	c := mqtt.NewClient(opts)

	// connect to broker
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Mqtt error: %s", token.Error())
	}

	// publicsh to broker
	for i := 0; i < 5; i++ {

		// to write file in pub server
		file_name_mnt := fmt.Sprintf("/mnt/test%d.text", i)
		d1 := []byte("hello world")
		err := os.WriteFile(file_name_mnt, d1, 0664)
		if err != nil {
			fmt.Println(err)
			return
		}

		// info of data
		nfs_server_addr := "10.0.8.19"
		nfs_server_port := "22"
		data_format := "file"
		file_name_nfs := fmt.Sprintf("/nfs/test%d.text", i)
		payload_data := Payload{
			Addr:     nfs_server_addr,
			Port:     nfs_server_port,
			Format:   data_format,
			Location: file_name_nfs}

		// to encode from golang structure to json
		jsonData, err := json.Marshal(payload_data)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("submit this data %s\n", jsonData)

		// publich to broker
		token := c.Publish("go-mqtt/sample", 0, false, jsonData)
		token.Wait()
	}

	c.Disconnect(250)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Complete publish")
}
