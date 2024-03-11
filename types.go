package dd_pubsub

import mqtt "github.com/eclipse/paho.mqtt.golang"

type Descriptor struct {
	Topic        string
	DataType     string
	Locator      string
	DatabaseAddr string
	DatabasePort string
	TimeStamp    string
	Header       string //Additional Data Information
}

type PubArg struct {
	Topic      string
	Qos        byte
	Retained   bool
	Payload    interface{}
	MqttClient mqtt.Client
	// Now only using object strage but it is ok that you can branch structured and unstructured data
	// RedisAddr  string //strage for structured_data
	// RedisPort  string //strage for structured_data
	StrageAddr string //strage for unstructured_data
	StragePort string //strage for unstructured_data
	StrageId   string
	StrageKey  string
}

type SubArg struct {
	Topic      string
	Qos        byte
	BrokerAddr string
	BrokerPort string
}
