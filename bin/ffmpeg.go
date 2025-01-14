package bin

import (
	"fmt"
	"io"
	"io/ioutil"
	"kowhai/global"
	"os"
	"os/exec"
)

func Start(hlsSegmentPattern, hlsM3U8File string, pr *io.PipeReader) error {
	useGPU := hasNvidiaGPU()
	var cmd *exec.Cmd

	// HLS 输出相关配置
	hlsSegmentTime := "6" // 每个 HLS 片段的时长（秒）

	if useGPU {
		cmd = exec.Command(
			"bin/ffmpeg",
			"-y",               // 覆盖输出文件
			"-hwaccel", "cuda", // 启用 GPU 加速
			"-i", "pipe:0", // 输入文件
			"-c:v", "h264_nvenc", // 使用 NVIDIA NVENC 编码器
			"-preset", "p5", // GPU 编码器预设
			"-hls_time", hlsSegmentTime, // 每个 HLS 片段的时长
			"-hls_segment_filename", hlsSegmentPattern, // ts 片段命名规则
			hlsM3U8File, // 输出 HLS 清单文件
		)
	} else {
		cmd = exec.Command(
			"bin/ffmpeg",
			"-y",           // 覆盖输出文件
			"-i", "pipe:0", // 输入文件
			"-c:v", "libx264", // 使用 CPU 编码器
			"-preset", "medium", // CPU 编码预设
			"-hls_time", hlsSegmentTime, // 每个 HLS 片段的时长
			"-hls_segment_filename", hlsSegmentPattern, // ts 片段命名规则
			hlsM3U8File, // 输出 HLS 清单文件
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
	return nil

}

// 获取视频时长
func GetVideoDuration(file io.Reader) (string, error) {
	// 创建临时文件
	tmpFile, err := ioutil.TempFile("", "upload_*.mp4")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer tmpFile.Close()

	// 将上传的文件内容写入临时文件
	if _, err := io.Copy(tmpFile, file); err != nil {
		return "", fmt.Errorf("写入临时文件失败: %w", err)
	}

	// 使用 ffprobe 获取视频信息
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", tmpFile.Name())
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("执行 ffprobe 失败: %w", err)
	}

	// 输出的是时长，去除尾部空格并返回
	duration := string(output)
	return duration, nil
}

func hasNvidiaGPU() bool {
	cmd := exec.Command("nvidia-smi")
	if err := cmd.Run(); err != nil {
		global.Logger.Error("No Nvidia GPU found", err)
		return false
	}
	return true
}
