package dd_pubsub

import (
	"encoding/json"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
	store "github.com/malcoscos/dd-pubsub/store_func"
	transport "github.com/malcoscos/dd-pubsub/transport_func"
	types "github.com/malcoscos/dd-pubsub/types"
)

func Publish(p *types.PubArg) {
	var data_mime_type string = transport.ProcessFile(p.Payload)
	var object_name string = uuid.NewString()

	if data_mime_type == "video" {
		object_name = store.StoreVideoData(p.Payload, object_name, p.MovieStrageDir)
	} else if data_mime_type == "image" || data_mime_type == "tiny_data" {
		object_name = store.StoreTinyData(p, object_name)
	}

	// descriptor of real data
	now := time.Now()
	time_stamp := fmt.Sprint(now.Format(time.RFC3339))
	descriptor := types.Descriptor{
		Topic:        p.Topic,
		DatabaseAddr: p.StrageAddr,
		DatabasePort: p.StragePort,
		DataType:     data_mime_type,
		Locator:      object_name,
		TimeStamp:    time_stamp,
		Header:       "hoge", // This attr is used after the ffmpeg implementation is finished
	}

	// to encode from golang structure to json
	jsonData, err := json.Marshal(descriptor)
	if err != nil {
		fmt.Println(err)
		return
	}

	// publich to broker
	token := p.MqttClient.Publish(p.Topic, p.Qos, p.Retained, jsonData)
	token.Wait()
	fmt.Println("Complete publish")
}
