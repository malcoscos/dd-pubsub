package dd_pubsub

import (
	"fmt"
	"net/http"
)

func ProcessFile(data interface{}) string {
	bytes, ok := data.([]byte)
	if !ok {
		return "unstructured_data"
	}

	mimeType := http.DetectContentType(bytes)

	switch {
	case mimeType == "video/mp4" || mimeType == "video/x-msvideo":
		fmt.Println("Processing video file with ffmpeg...")
		// ここでffmpegを実行します。
		// 実際にはffmpegコマンドはファイルを要求するため、ファイルへの書き出しが必要です。
		return "unstructured_data"
	case mimeType == "image/jpeg" || mimeType == "image/png":
		fmt.Println("Processing image file with jhead...")
		// ここでjheadを実行します。
		// 実際にはjheadコマンドもファイルを要求するため、ファイルへの書き出しが必要です。
		return "unstructured_data"
	default:
		fmt.Println("No action required for this file type.")
		return "structured_data"
	}
}
