package bin

import (
	"io"
	"kowhai/apps/minio"
	"kowhai/global"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Start(ts, m3u8, minio_path string, hlsDir string, userId int, pr *io.PipeReader) error {
	useGPU := hasNvidiaGPU()
	var cmd *exec.Cmd

	// HLS 输出相关配置
	hlsSegmentTime := "10" // 每个 HLS 片段的时长（秒）

	tempDir := "/tmp"
	// TS 和M3U8临时文件路径
	tsFilePathPattern := filepath.Join(tempDir, ts)
	m3u8FilePath := filepath.Join(tempDir, m3u8)

	if useGPU {
		cmd = exec.Command(
			"bin/ffmpeg",
			"-y",               // 覆盖输出文件
			"-hwaccel", "cuda", // 启用 GPU 加速
			"-i", "pipe:0", // 输入文件
			"-c:v", "h264_nvenc", // 使用 NVIDIA NVENC 编码器
			"-preset", "p5", // GPU 编码器预设
			"-hls_time", hlsSegmentTime, // 每个 HLS 片段的时长
			"-hls_playlist_type", "vod",
			"-hls_segment_filename", tsFilePathPattern, // ts 片段命名规则
			"-hls_base_url", minio_path,
			m3u8FilePath, // 输出 HLS 清单文件
		)
	} else {
		cmd = exec.Command(
			"bin/ffmpeg",
			"-y",           // 覆盖输出文件
			"-i", "pipe:0", // 输入文件
			"-c:v", "libx264", // 使用 CPU 编码器
			"-preset", "medium", // CPU 编码预设
			"-hls_time", hlsSegmentTime, // 每个 HLS 片段的时长
			"-hls_playlist_type", "vod",
			"-hls_segment_filename", tsFilePathPattern, // ts 片段命名规则
			"-hls_base_url", minio_path,
			m3u8FilePath, // 输出 HLS 清单文件
		)
	}

	cmd.Stdin = pr

	cmd.Stdout = os.Stdout

	cmd.Stderr = os.Stderr
	// 打印 FFmpeg 命令
	global.Logger.Info("Start ffmpeg command", "cmd", cmd.String())

	// 启动 FFmpeg 命令
	if err := cmd.Start(); err != nil {
		global.Logger.Error("Failed to start ffmpeg command", err)
		return err
	}

	// 等待 FFmpeg 处理完成
	if err := cmd.Wait(); err != nil {
		global.Logger.Error("Failed to wait ffmpeg command", err)
		return err
	}

	// 上传逻辑
	tsPattern := strings.Replace(ts, "%03d", "*", 1) // 将 %03d 替换为 *，变成 faruxue_*.ts
	pattern := filepath.Join(tempDir, tsPattern)     // 生成匹配模式

	global.Logger.Info("Glob pattern", "pattern", pattern)

	// 使用 glob 来查找所有匹配的文件
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

	// 上传 M3U8 文件
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
