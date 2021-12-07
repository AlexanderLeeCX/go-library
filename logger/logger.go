/**
 * @Author: Lee
 * @Description:
 * @File:  logger
 * @Version: 1.0.0
 * @Date: 2021/12/7 10:06 下午
 */

package logger

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

// 日志配置结构体
type LogEncoderConfig struct {
	zapcore.EncoderConfig
	lumberjack.Logger
	LoggerLevel   string            // 日志等级
	OutputConsole bool              // 是否输出到控制台
	OutputFile    bool              // 是否输出到文件
	FormatJson    bool              // 是否输出json格式
	FixedFieldMap map[string]string // 固定字段
	isGRPC        bool              // 是否为grpc服务
}

// 日志结构体
type Logger struct {
	*zap.Logger
	logConfig *LogEncoderConfig
}

// 获取日志等级对象
func getLogLevel(levelStr string) (level zapcore.Level, err error) {
	switch levelStr {
	case "debug":
		level = zap.DebugLevel
		break
	case "info":
		level = zap.InfoLevel
		break
	case "warn":
		level = zap.WarnLevel
		break
	case "error":
		level = zap.ErrorLevel
		break
	case "dpanic":
		level = zap.DPanicLevel
		break
	case "panic":
		level = zap.PanicLevel
		break
	case "fatal":
		level = zap.FatalLevel
		break
	default:
		return -2, errors.New("log level not found!")
	}
	return
}

// InitLogger 初始化日志配置
func InitLogger(config *LogEncoderConfig) (*Logger, error) {

	loggerLevelStr := strings.ToLower(config.LoggerLevel)
	loggerLevel, err := getLogLevel(loggerLevelStr)
	if err != nil {
		return nil, err
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(loggerLevel)

	// 判断是否输出控制台或者文件
	fileSyncList := make([]zapcore.WriteSyncer, 0)
	consoleSyncList := make([]zapcore.WriteSyncer, 0)
	if config.OutputConsole {
		consoleSyncList = append(consoleSyncList, zapcore.AddSync(os.Stdout))
	}
	if config.OutputFile {
		fileSyncList = append(fileSyncList, zapcore.AddSync(&config.Logger))
	}

	// 判断是否输出json格式
	var fileEncoder zapcore.Encoder
	consoleEncoder := zapcore.NewConsoleEncoder(config.EncoderConfig)
	if config.FormatJson {
		fileEncoder = zapcore.NewJSONEncoder(config.EncoderConfig)
	} else {

	}

	fileCore := zapcore.NewCore(
		fileEncoder, // 编码器配置,
		zapcore.NewMultiWriteSyncer(fileSyncList...), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	consoleCore := zapcore.NewCore(
		consoleEncoder, // 编码器配置,
		zapcore.NewMultiWriteSyncer(consoleSyncList...), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	core := zapcore.NewTee(fileCore, consoleCore)

	optionList := make([]zap.Option, 0)

	// 堆栈跟踪
	caller := zap.AddCaller()
	optionList = append(optionList, caller)

	// 设置初始化字段
	if len(config.FixedFieldMap) > 0 {
		fields := make([]zap.Field, 0)
		for key, value := range config.FixedFieldMap {
			fields = append(fields, zap.String(key, value))
		}
		optionList = append(optionList, zap.Fields(fields...))
	}

	logger := &Logger{Logger: zap.New(core, optionList...)}
	logger.logConfig = config
	return logger, nil
}

// 获取额外的日志字段
func _getEnvFields() []zap.Field {
	fields := make([]zap.Field, 0)

	// 日志追加容器信息
	envMaps := make([]map[string]string, 0)
	for _, envMap := range envMaps {
		fields = append(fields, zapcore.Field{
			Type:   zapcore.StringType,
			Key:    envMap["key"],
			String: envMap["value"],
		})
	}

	return fields
}

// grpc打印DEBUG日志
func (logger *Logger) Debug(msg string) {
	fields := _getEnvFields()
	logger.Logger.Debug(msg, fields...)
}

// grpc打印ERROR日志, msg传入error.Error()信息
func (logger *Logger) Error(err error) {
	fields := _getEnvFields()
	stackTraceField := zapcore.Field{
		Type:   zapcore.StringType,
		Key:    "stacktrace",
		String: fmt.Sprintf("%+v", errors.WithStack(err)),
	}
	fields = append(fields, stackTraceField)
	logger.Logger.Error(err.Error(), fields...)
	fmt.Println(fmt.Sprintf("%+v", errors.WithStack(err)))
}

// grpc打印INFO日志
func (logger *Logger) Info(msg string) {
	fields := _getEnvFields()
	logger.Logger.Info(msg, fields...)
}

// grpc打印Warn日志
func (logger *Logger) Warn(msg string) {
	fields := _getEnvFields()
	logger.Logger.Warn(msg, fields...)
}

// grpc打印DPanic日志
func (logger *Logger) DPanic(msg string) {
	fields := _getEnvFields()
	logger.Logger.DPanic(msg, fields...)
}

// grpc打印Fatal日志
func (logger *Logger) Fatal(msg string) {
	fields := _getEnvFields()
	logger.Logger.Fatal(msg, fields...)
}

// grpc打印Debugf日志
func (logger Logger) Debugf(template string, args ...interface{}) {
	fields := _getEnvFields()
	logger.Logger.Debug(fmt.Sprintf(template, args...), fields...)
}

// grpc打印Infof日志
func (logger Logger) Infof(template string, args ...interface{}) {
	fields := _getEnvFields()
	logger.Logger.Info(fmt.Sprintf(template, args...), fields...)
}

// grpc打印Warnf日志
func (logger Logger) Warnf(template string, args ...interface{}) {
	fields := _getEnvFields()
	logger.Logger.Warn(fmt.Sprintf(template, args...), fields...)
}

// grpc打印Errorf日志
func (logger Logger) Errorf(template string, err error, args ...interface{}) {
	fields := _getEnvFields()
	stackTraceField := zapcore.Field{
		Type:   zapcore.StringType,
		Key:    "stacktrace",
		String: fmt.Sprintf("%+v", errors.WithStack(err)),
	}
	fields = append(fields, stackTraceField)
	logger.Logger.Error(fmt.Sprintf(template, args...), fields...)
	fmt.Println(fmt.Sprintf("%+v", errors.WithStack(err)))
}

// grpc打印DPanicf日志
func (logger Logger) DPanicf(template string, args ...interface{}) {
	fields := _getEnvFields()
	logger.Logger.DPanic(fmt.Sprintf(template, args...), fields...)
}

// grpc打印Panicf日志
func (logger Logger) Panicf(template string, args ...interface{}) {
	fields := _getEnvFields()
	logger.Logger.Panic(fmt.Sprintf(template, args...), fields...)
}

// grpc打印Fatalf日志
func (logger Logger) Fatalf(template string, args ...interface{}) {
	fields := _getEnvFields()
	logger.Logger.Fatal(fmt.Sprintf(template, args...), fields...)
}
