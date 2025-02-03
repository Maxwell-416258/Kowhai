package bin

import (
	"github.com/u2takey/ffmpeg-go"
	"io"
	"kowhai/apps/streaming/minio"
	"kowhai/global"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Start(ts, m3u8, minioPath string, hlsDir string, userId int, pr *io.PipeReader) error {
	useGPU := hasNvidiaGPU()
	hlsSegmentTime := "6" // HLS 片段时长（秒）
	tempDir := "/tmp"
	tsFilePathPattern := filepath.Join(tempDir, ts)
	m3u8FilePath := filepath.Join(tempDir, m3u8)

	var ffmpegCmd *ffmpeg_go.Stream

	if useGPU {
		ffmpegCmd = ffmpeg_go.Input("pipe:0",
			ffmpeg_go.KwArgs{
				"hwaccel":               "cuda",
				"hwaccel_output_format": "cuda",
			}).
			Output(m3u8FilePath,
				ffmpeg_go.KwArgs{
					"c:v":                  "h264_nvenc", // 使用 GPU
					"preset":               "p6",
					"hls_time":             hlsSegmentTime,
					"hls_playlist_type":    "vod",
					"hls_segment_filename": tsFilePathPattern,
					"hls_base_url":         minioPath,
				})
	} else {
		ffmpegCmd = ffmpeg_go.Input("pipe:0").
			Output(m3u8FilePath,
				ffmpeg_go.KwArgs{
					"c:v":                  "libx264", // 使用 CPU
					"preset":               "medium",
					"hls_time":             hlsSegmentTime,
					"hls_playlist_type":    "vod",
					"hls_segment_filename": tsFilePathPattern,
					"hls_base_url":         minioPath,
				})
	}

	ffmpegCmd = ffmpegCmd.WithInput(pr)
	global.Logger.Info("Start ffmpeg command")

	if err := ffmpegCmd.Run(); err != nil {
		global.Logger.Error("Failed to run ffmpeg", err)
		return err
	}

	tsPattern := strings.Replace(ts, "%03d", "*", 1)
	pattern := filepath.Join(tempDir, tsPattern)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		global.Logger.Error("Failed to find TS files", err)
		return err
	}

	for _, tsFilePath := range matches {
		tsFileName := filepath.Base(tsFilePath)
		tsFile, err := os.Open(tsFilePath)
		if err != nil {
			global.Logger.Error("Failed to open TS file", err, "file", tsFileName)
			continue
		}

		uploadErr := minio.UploadVideo(userId, hlsDir, tsFileName, tsFile)
		tsFile.Close()
		if uploadErr != nil {
			global.Logger.Error("Failed to upload TS file", uploadErr, "file", tsFileName)
		}
	}

	go func() {
		m3u8File, err := os.Open(m3u8FilePath)
		if err != nil {
			global.Logger.Error("Failed to open M3U8 file", err)
			return
		}
		defer m3u8File.Close()

		uploadErr := minio.UploadVideo(userId, hlsDir, m3u8, m3u8File)
		if uploadErr != nil {
			global.Logger.Error("Failed to upload M3U8 file", uploadErr)
		}
	}()

	return nil
}

func hasNvidiaGPU() bool {
	cmd := exec.Command("nvidia-smi")
	if err := cmd.Run(); err != nil {
		global.Logger.Error("No Nvidia GPU found", err)
		return false
	}
	return true
}
