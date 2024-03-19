package store_func

import (
	"fmt"
	"os"
	"path/filepath"
)

func StoreVideoData(data interface{}, object_name string, dir string) string {
	payload_data, ok := data.([]byte)
	if !ok {
		fmt.Println("Failed to exchange data to byte", ok)
		return ""
	}
	// データをファイルに保存し、ファイルパスを取得
	filepath, err := SaveDataToFile(payload_data, dir, object_name)
	if err != nil {
		fmt.Println("Failed to save data to file:", err)
		return ""
	}
	fmt.Println("Data saved to file:", filepath)
	return filepath
}

func SaveDataToFile(data []byte, dir, file_name string) (string, error) {
	// ディレクトリを作成（存在しない場合）
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// ファイル名を生成（プレフィックス+タイムスタンプ）
	fullPath := filepath.Join(dir, file_name)

	// ファイルを開く（存在しない場合は作成、存在する場合は上書き）
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// データをファイルに書き込む
	if _, err = file.Write(data); err != nil {
		return "", err
	}

	return fullPath, nil
}
