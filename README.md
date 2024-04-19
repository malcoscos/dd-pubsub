# DD-PubSub
## 目次
- システム構成概要
- 使用技術
- 環境構築
- ディレクトリ構造

## システム構成概要
![alt text](image-1.png)
## システムフロー
![alt text](image-2.png)
## Prerequisites
- Go version 1.15 or higher
- Docker 
## 環境構築
### Install DD-PubSub using go get:
``` bash
go get github.com/malcoscos/dd-pubsub
go get github.com/eclipse/paho.mqtt.golang
```
### Install mosquitto
brokerとしてmosquittoをインストール
``` bash
docker pull eclipse-mosquitto
docker run -it -p 1883:1883  --name mosquitto eclipse-mosquitto
```
### MinIOストレージの構築
ストレージ用のサーバーに対してMinIOストレージを構築
``` bash
docker pull minio/minio
docker run -p 9000:9000 --name minio1 -e "MINIO_ROOT_USER=youraccesskey" -e "MINIO_ROOT_PASSWORD=yoursecretkey" -v /mnt/data:/data minio/minio server /data --console-address ":9001"
```
## Usage
### Publisher
``` golang
import (
	"fmt"
	"log"
	"net/http"
	"os"

	dd_pubsub "github.com/malcoscos/dd-pubsub"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
  // configure mqtt client options
  opts := mqtt.NewClientOptions()
  // add broker
  broker := fmt.Sprintf("tcp://%s:%s", "127.0.0.1", "1883")
  opts.AddBroker(broker)
  // create mqtt client
  c := mqtt.NewClient(opts)
  // publish data using dd-pubsub library
  pub_arg := dd_pubsub.PubArg{
  	Topic:      "hoge", // topic of publishing data
  	Qos:        0,  // Qos level of publish message
  	Retained:   false, // retain message in broker
  	Payload:    "hoge",
  	MqttClient: c,
  	StrageAddr: "127.0.0.1", // strage addr
  	StragePort: "9000", // strage port num
  	StrageID:   "hoge", // strage id
  	StrageKey:  "hoge", // strage key of using 
  }

  dd_pubsub.Publish(&pub_arg)
}
```
### Subscriber
```golang
import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	dd_pubsub "github.com/malcoscos/dd-pubsub"
)

func main(){
  // configure mqtt client options
	opts := mqtt.NewClientOptions()
	// add broker
	broker := fmt.Sprintf("tcp://%s:%s", "127.0.0.1", "1883")
	opts.AddBroker(broker)
	// create mqtt client
	c := mqtt.NewClient(opts)
	// subscribe data using dd-pubsub library
	sub_arg := dd_pubsub.SubArg{
		Topic:      "demo", // subscribe topic name
		Qos:        0, // Qos level of subscribe message
		MqttClient: c,
		StorePath:  "/path/to/subscribe/data",
	}

	dd_pubsub.Subscribe(&sub_arg)
}
```


## ディレクトリ構造