package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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
		d1 := []byte("hello world")
		file_name := fmt.Sprintf("/tmp/test%d.text", i)
		err := os.WriteFile(file_name, d1, 0664)
		if err != nil {
			fmt.Println(err)
			return
		}
		// text := fmt.Sprintf("this is msg #%d!", i)

		token := c.Publish("go-mqtt/sample", 0, false, file_name)

		token.Wait()
	}

	c.Disconnect(250)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Complete publish")
}

type data struct {
	pub_data int
}
