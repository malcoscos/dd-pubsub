package retrieve_func

import (
	"fmt"
	"os"

	"github.com/gorilla/websocket"
	"github.com/malcoscos/dd-pubsub/types"
)

func RetreiveVideoData(payload_data string, descriptor types.Descriptor, s *types.SubArg) {
	// connect to server of websocket
	url := fmt.Sprintf("ws://%s/%s", descriptor.DatabaseAddr, descriptor.Locator)
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	// サーバーからのメッセージを受信して表示するループ
	_, message, err := conn.ReadMessage()
	if err != nil {
		fmt.Println(err)
	}
	// 受信したメッセージを表示
	fmt.Printf("受信メッセージ: %s\n", message)

	// create file
	file_path := fmt.Sprintf("%s/%s", s.StorePath, descriptor.TimeStamp)
	file, err := os.Create(file_path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// write file
	_, err = file.Write(message)
	if err != nil {
		fmt.Println(err)
	}
}
