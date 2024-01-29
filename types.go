package dd_pubsub

type Descriptor struct {
	Topic        string
	DataType     string
	Locator      string
	DatabaseAddr string
	DatabasePort string
	TimeStamp    string
	Header       string
}

type PubArg struct {
	Topic      string
	Qos        byte
	Retained   bool
	Payload    interface{}
	BrokerAddr string
	BrokerPort string
	RedisAddr  string
	RedisPort  string
	MinioAddr  string
	MinioPort  string
}

type SubArg struct {
	Topic      string
	Qos        byte
	BrokerAddr string
	BrokerPort string
}
