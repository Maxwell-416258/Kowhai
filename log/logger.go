package log

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"time"
	"vidspark/global"
)

func InitLogger() *zap.SugaredLogger {
	logMode := zapcore.DebugLevel
	if global.Config.Log.Develop {
		logMode = zapcore.InfoLevel
	}
	core := zapcore.NewCore(getEncoder(), getWriteSyncer(), logMode)
	return zap.New(core).Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 日志级别改大写
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Local().Format(time.DateTime)) //time格式化
	}
	return zapcore.NewJSONEncoder(encoderConfig)

}

// 定义日志输出地方
func getWriteSyncer() zapcore.WriteSyncer {
	stSeparator := string(filepath.Separator)
	stRootDir, _ := os.Getwd()
	stLogFilePath := stRootDir + stSeparator + "log" + stSeparator + "logg" + stSeparator + time.Now().Format(time.DateOnly) + ".log"

	fmt.Println(stLogFilePath)

	luberjackSyncer := &lumberjack.Logger{
		Filename:   stLogFilePath,
		MaxSize:    global.Config.Log.MaxSize,
		MaxBackups: global.Config.Log.MaxBackups,
		MaxAge:     global.Config.Log.MaxAge,
		Compress:   global.Config.Log.Compress,
	}
	return zapcore.AddSync(luberjackSyncer)
}
