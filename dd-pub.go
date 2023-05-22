package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// json formatの作成
type Payload struct {
	Addr    string
	Port    int
	Format  string
	Locator string
}

func main() {

	// optsにClientOptionsインスタンスのpointerを格納
	opts := mqtt.NewClientOptions()

	//　BrokerServerのlistに追加
	opts.AddBroker("tcp://10.0.8.25:1883")

	// clientクラスのインスタンスを作成
	c := mqtt.NewClient(opts)
	// BrokerへのconnectionにErrorがないか判定
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Mqtt error: %s", token.Error())
	}
	// clientからosのシステムを利用してパケットをbrokerにstoreする
	for i := 0; i < 5; i++ {
		nfs_server_addr := "10.0.8.19"
		nfs_server_port := 22
		data_format := "file"

		d1 := []byte("hello world")
		file_name_mnt := fmt.Sprintf("/mnt/test%d.text", i)

		err := os.WriteFile(file_name_mnt, d1, 0664)
		if err != nil {
			fmt.Println(err)
			return
		}

		file_name_nfs := fmt.Sprintf("/nfs/test%d.text", i)

		payload_data := Payload{
			Addr:    nfs_server_addr,
			Port:    nfs_server_port,
			Format:  data_format,
			Locator: file_name_nfs}

		jsonData, err := json.Marshal(payload_data)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%s\n", jsonData)

		token := c.Publish("go-mqtt/sample", 0, false, jsonData)

		token.Wait()
	}

	c.Disconnect(250)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Complete publish")
}
