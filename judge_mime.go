package dd_pubsub

import (
	"fmt"
	"net/http"
)

func ProcessFile(data interface{}) string {
	bytes, ok := data.([]byte)
	if !ok {
		return "tiny_data"
	}

	mimeType := http.DetectContentType(bytes)

	switch {
	case mimeType == "video/mp4" || mimeType == "video/x-msvideo":

		// fmt.Println("Processing video file with ffmpeg...")
		// // 変換するビデオファイルのパス
		// inputFile := "input.mp4"
		// outputFile := "output.avi"

		// // FFmpegコマンドを構築
		// cmd := exec.Command("ffmpeg", "-i", inputFile, outputFile)

		// // コマンドを実行
		// if err := cmd.Run(); err != nil {
		// fmt.Println("Error executing FFmpeg command:", err)
		// return
		// }
		// fmt.Println("Conversion complete, file saved as:", outputFile)

		return "video_data"
	case mimeType == "image/jpeg" || mimeType == "image/png":
		fmt.Println("Processing image file with jhead...")
		// ここでjheadを実行します。
		// 実際にはjheadコマンドもファイルを要求するため、ファイルへの書き出しが必要です。
		return "image_data"
	default:
		fmt.Println("No action required for this file type.")
		return "tiny_data"
	}
}
