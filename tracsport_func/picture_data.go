package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net"

	"github.com/brocaar/lorawan"
)

type Data struct {
	Type string
	Info string
}

func main() {
	// Read image file
	imageBytes, err := ioutil.ReadFile("image.jpg")
	if err != nil {
		panic(err)
	}

	// Convert image data to base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	// Create data object
	data := Data{
		Type: "Image",
		Info: imageBase64,
	}

	// Convert data to JSON
	dataBytes, _ := json.Marshal(data)

	// Initialize variables
	var lorawanMACPayload lorawan.MACPayload
	var lorawanPHYPayload lorawan.PHYPayload

	// Convert data to MACPayload
	lorawanMACPayload.FPort = new(uint8) // Set FPort value to 1
	lorawanMACPayload.FRMPayload = []lorawan.Payload{
		&lorawan.DataPayload{Bytes: dataBytes},
	}

	// Create a PHYPayload from the MACPayload
	lorawanPHYPayload.MHDR = lorawan.MHDR{MType: lorawan.UnconfirmedDataUp}
	lorawanPHYPayload.MACPayload = &lorawanMACPayload

	// Your LoRaWAN session keys
	var appSKey lorawan.AES128Key // Application Session Key

	// Encrypt the payload
	err = lorawanPHYPayload.EncryptFRMPayload(appSKey)
	if err != nil {
		panic(err)
	}

	// Get the encoded bytes
	phyBytes, err := lorawanPHYPayload.MarshalBinary()
	if err != nil {
		panic(err)
	}

	// Send data over socket
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	_, err = conn.Write(phyBytes)
	if err != nil {
		panic(err)
	}
}
