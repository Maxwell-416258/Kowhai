package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewLogger(filePath string) *zap.Logger {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// Create a write syncer
	writerSyncer := zapcore.AddSync(file)

	// Define the encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "message",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	// Create a core with the file write syncer and JSON encoding
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writerSyncer,
		zapcore.InfoLevel, // 最低日志输出级别
	)
	return zap.New(core)
}

var (
	SQLLogger    *zap.Logger
	SystemLogger *zap.Logger
)

func InitLoggers() {
	SQLLogger = NewLogger("sql.log")
	SystemLogger = NewLogger("system.log")
}
