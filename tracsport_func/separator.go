package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/brocaar/lorawan"
)

type Data struct {
	Type string
	Info string
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}

		// get data from client
		buffer, err := ioutil.ReadAll(conn)
		if err != nil {
			panic(err)
		}

		// initialize variable
		var lorawanPHYPayload lorawan.PHYPayload

		// create a PHYPayload from the received bytes
		err = lorawanPHYPayload.UnmarshalBinary(buffer)
		if err != nil {
			panic(err)
		}

		// your LoRaWAN session keys
		var appSKey lorawan.AES128Key // Application Session Key

		// decrypt the payload
		err = lorawanPHYPayload.DecryptFRMPayload(appSKey)
		if err != nil {
			panic(err)
		}

		// interpret data
		if len(lorawanPHYPayload.MACPayload.(*lorawan.MACPayload).FRMPayload) > 0 {
			payload := lorawanPHYPayload.MACPayload.(*lorawan.MACPayload).FRMPayload[0]

			var data Data
			err = json.Unmarshal(payload.(*lorawan.DataPayload).Bytes, &data)
			if err != nil {
				panic(err)
			}

			// handle data based on type
			switch data.Type {
			case "Temperature":
				fmt.Println("Temperature Data:", data.Info)
			case "Pressure":
				fmt.Println("Pressure Data:", data.Info)
			case "Image":
				fmt.Println("Image Data:", data.Info)
			case "Video":
				fmt.Println("Video Data:", data.Info)
			}
		}
	}
}
