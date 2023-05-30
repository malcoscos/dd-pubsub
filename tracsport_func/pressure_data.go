package main

import (
	"encoding/json"
	"net"

	lorawan "github.com/brocaar/lorawan"
)

type Data struct {
	Type string
	Info string
}

func main() {
	data := Data{"Pressure", "27.3"}

	dataBytes, _ := json.Marshal(data)

	// initialize variables
	var lorawanMACPayload lorawan.MACPayload
	var lorawanPHYPayload lorawan.PHYPayload

	// convert the data to MACPayload
	fPort := uint8(1) // Create a temporary variable and assign the value 1
	lorawanMACPayload.FPort = &fPort
	lorawanMACPayload.FRMPayload = []lorawan.Payload{
		&lorawan.DataPayload{Bytes: dataBytes},
	}

	// create a PHYPayload from the MACPayload
	lorawanPHYPayload.MHDR = lorawan.MHDR{MType: lorawan.UnconfirmedDataUp}
	lorawanPHYPayload.MACPayload = &lorawanMACPayload

	// your LoRaWAN session keys
	var appSKey lorawan.AES128Key // Application Session Key

	// encrypt the payload
	err := lorawanPHYPayload.EncryptFRMPayload(appSKey)
	if err != nil {
		panic(err)
	}

	// get the encoded bytes
	phyBytes, err := lorawanPHYPayload.MarshalBinary()
	if err != nil {
		panic(err)
	}

	// send data over socket
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	_, err = conn.Write(phyBytes)
	if err != nil {
		panic(err)
	}
}
