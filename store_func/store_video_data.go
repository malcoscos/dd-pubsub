package store_func

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func StoreVideoData(data interface{}, object_name string, dir string) string {
	payload_data, ok := data.([]byte)
	if !ok {
		fmt.Println("Failed to exchange data to byte", ok)
		return ""
	}
	// Save the data to a file and obtain the file path
	filepath, err := SaveDataToFile(payload_data, dir, object_name)
	if err != nil {
		fmt.Println("Failed to save data to file:", err)
		return ""
	}
	fmt.Println("Data saved to file:", filepath)
	return filepath
}

func SaveDataToFile(data []byte, dir, file_name string) (string, error) {
	// Create directory if it does not exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// Generate file name (prefix + timestamp)
	fullPath := filepath.Join(dir, file_name)

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "example_*.mp4")
	if err != nil {
		fmt.Printf("Failed to create temp file: %v\n", err)
		return "", err
	}
	defer os.Remove(tempFile.Name()) // Remove the temporary file after completion

	// Write data to the temporary file (as an example, write "Hello, world!")
	if _, err := tempFile.Write(data); err != nil {
		fmt.Printf("Failed to write to temp file: %v\n", err)
		return "", err
	}

	// Close the temporary file
	if err := tempFile.Close(); err != nil {
		fmt.Printf("Failed to close temp file: %v\n", err)
		return "", err
	}

	// Compress the file using ffmpeg
	cmd := exec.Command("ffmpeg", "-i", tempFile.Name(), "-b:v", "1M", "-c:a", "copy", fullPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Failed to compress video: %v, ffmpeg output: %s\n", err, output)
		return "", err
	}

	return fullPath, nil
}
