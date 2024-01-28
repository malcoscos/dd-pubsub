package dd_pubsub

type Descriptor struct {
	Format       string
	Locator      string
	DatabaseAddr string
	DatabasePort string
	TimeStamp    string
	Header       string
}

type PubArg struct {
	Topic        string
	Qos          byte
	Retained     bool
	Payload      interface{}
	DataFormat   string
	BrokerAddr   string
	BrokerPort   string
	DatabaseAddr string
	DatabasePort string
}

type SubArg struct {
	Topic      string
	Qos        byte
	BrokerAddr string
	BrokerPort string
}
