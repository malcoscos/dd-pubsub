package dd_pubsub

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	retrieve_func "github.com/malcoscos/dd-pubsub/retrieve_func"
	types "github.com/malcoscos/dd-pubsub/types"
)

func Subscribe(s *types.SubArg) {

	// make channel
	msgCh := make(chan mqtt.Message)
	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		msgCh <- msg
	}
	// mqtt client
	c := s.MqttClient
	//connect to broker
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	// subscribe from broker
	if subscribeToken := c.Subscribe(s.Topic, s.Qos, f); subscribeToken.Wait() && subscribeToken.Error() != nil {
		fmt.Println(subscribeToken.Error())
	}

	//  notify systemcall get from channel
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	// forever
	for {
		select {
		// get message from channel
		case m := <-msgCh:
			var descriptor types.Descriptor
			payload_data := string(m.Payload())

			// to decode from golong structure to json
			if err := json.Unmarshal([]byte(payload_data), &descriptor); err != nil {
				fmt.Println(err)
				return
			}

			if descriptor.DataType == "video_data" {
				retrieve_func.RetreiveVideoData(payload_data, descriptor, s)
			} else if descriptor.DataType == "image" || descriptor.DataType == "tiny_data" {
				retrieve_func.RetreiveTinyData(descriptor, s)
			}
		// to interrupt if there is systemcall
		case <-signalCh:
			fmt.Printf("Interrupt detected.\n")
			c.Disconnect(1000)
			return
		}
	}
}
