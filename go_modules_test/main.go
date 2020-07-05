package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)


func getLogWriter()zapcore.WriteSyncer {
	file, _ := os.Create("./test.log")
	return zapcore.AddSync(file)
}

func Init() {
	writeSyncer := getLogWriter() // 日志输出位置相关
	encoder := getEncoder() // 日志的格式相关
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	logger = logger
	//sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 人类可读的时间格式
	//encoderConfig.EncodeTime = zapcore.EpochTimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	//return zapcore.NewConsoleEncoder(encoderConfig)  // 可读日志
	return zapcore.NewJSONEncoder(encoderConfig) // json格式日志
}


func main() {



	fmt.Println("OK")




}
