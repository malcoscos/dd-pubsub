package transport_func

import (
	"fmt"
	"net/http"
)

func ProcessData(data interface{}) string {
	bytes, ok := data.([]byte)
	if !ok {
		return "tiny_data"
	}

	mimeType := http.DetectContentType(bytes)

	switch {
	case mimeType == "video/mp4" || mimeType == "video/x-msvideo":
		return "video_data"
	case mimeType == "image/jpeg" || mimeType == "image/png":
		fmt.Println("Processing image file with jhead...")
		return "image_data"
	default:
		fmt.Println("No action required for this file type.")
		return "tiny_data"
	}
}
