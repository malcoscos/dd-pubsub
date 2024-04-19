// package transport_func

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"os/exec"
// )

// func ExcuteFFMPEG(movie_data interface{}) {
// 	// 一時ファイルを作成して動画データを書き込む
// 	tmpFile, err := ioutil.TempFile("", "video-*.mp4")
// 	if err != nil {
// 		return fmt.Errorf("failed to create temp file: %s", err)
// 	}
// 	defer os.Remove(tmpFile.Name()) // 関数終了時に一時ファイルを削除

// 	if _, err := tmpFile.Write(videoData); err != nil {
// 		return fmt.Errorf("failed to write video data to temp file: %s", err)
// 	}
// 	if err := tmpFile.Close(); err != nil {
// 		return fmt.Errorf("failed to close temp file: %s", err)
// 	}

// 	// ffmpegコマンドを実行して動画を処理
// 	// ここでは例として、動画のフォーマット情報を取得して表示します
// 	cmd := exec.Command("ffmpeg", "-i", tmpFile.Name(), "-f", "ffmetadata")
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("ffmpeg command failed: %s\n%s", err, string(output))
// 	}

// 	fmt.Printf("ffmpeg output: \n%s\n", string(output))
// 	return nil
// }

// func ExcuteJhead(image_data interface{}) {
// 	// 一時ファイルを作成
// 	tmpFile, err := os.TempFile("", "image-*.jpg")
// 	if err != nil {
// 		return fmt.Errorf("failed to create temp file: %s", err)
// 	}
// 	defer os.Remove(tmpFile.Name()) // 関数終了時に一時ファイルを削除

// 	// 画像データを一時ファイルに書き込む
// 	if _, err := tmpFile.Write(imageData); err != nil {
// 		return fmt.Errorf("failed to write image data to temp file: %s", err)
// 	}
// 	if err := tmpFile.Close(); err != nil {
// 		return fmt.Errorf("failed to close temp file: %s", err)
// 	}

// 	// jheadを実行して画像ファイルを処理
// 	cmd := exec.Command("jhead", tmpFile.Name())
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("jhead command failed: %s\n%s", err, string(output))
// 	}

// 	fmt.Printf("jhead output: \n%s\n", string(output))
// 	return nil
// }
