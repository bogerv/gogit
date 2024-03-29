// Package logx
// 参照 Tony Bai 的文章实现
// https://tonybai.com/2021/07/14/uber-zap-advanced-usage/
package logx

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

const (
	MaxSize    = 50 // 每个日志文件最大尺寸(M)
	MaxBackups = 10 // 日志文件最多保存20个备份
	MaxAge     = 15 // 天
)

type RotateOptions struct {
	MaxSize    int  // 每个日志文件保存的最大尺寸 单位：M
	MaxAge     int  // 文件最多保存多少天
	MaxBackups int  // 日志文件最多保存多少个备份
	Compress   bool // 是否压缩
}

type TeeOption struct {
	Filename     string // 日志文件路径
	Rotate       RotateOptions
	LevelEnabler zap.LevelEnablerFunc
}

// Init 初始化全局日志实例
func Init(filename string) {
	paths, _ := filepath.Split(filename)
	base := path.Base(filename)
	suffix := path.Ext(filename)
	prefix := strings.TrimSuffix(base, suffix)

	access := path.Join(paths, prefix+".access"+suffix)
	err := path.Join(paths, prefix+".err"+suffix)
	var tops = []TeeOption{
		{
			Filename: access,
			Rotate: RotateOptions{
				MaxSize:    MaxSize,
				MaxAge:     MaxAge,
				MaxBackups: MaxBackups,
				Compress:   true,
			},
			LevelEnabler: func(lvl zapcore.Level) bool {
				return lvl <= zap.InfoLevel
			},
		},
		{
			Filename: err,
			Rotate: RotateOptions{
				MaxSize:    MaxSize,
				MaxAge:     MaxAge,
				MaxBackups: MaxBackups,
				Compress:   true,
			},
			LevelEnabler: func(lvl zapcore.Level) bool {
				return lvl > zap.InfoLevel
			},
		},
	}

	logger = New(tops, zap.AddCaller(), zap.AddCallerSkip(1))
}

// New 新创建实例化日志
func New(tops []TeeOption, opts ...zap.Option) *zap.Logger {
	var cores []zapcore.Core
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.CallerKey = "source"
	//cfg.EncoderConfig.EncodeName = zapcore.FullNameEncoder
	// 定义时间格式
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	}

	for _, top := range tops {
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   top.Filename,          // 日志文件路径
			MaxSize:    top.Rotate.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: top.Rotate.MaxBackups, // 日志文件最多保存多少个备份
			MaxAge:     top.Rotate.MaxAge,     // 文件最多保存多少天
			Compress:   top.Rotate.Compress,   // 是否压缩
			LocalTime:  true,                  // 备份文件名本地/UTC时间
		})
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),                  // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), w), // 打印到控制台和文件
			top.LevelEnabler, // 日志级别
		)
		cores = append(cores, core)
	}

	return zap.New(zapcore.NewTee(cores...), opts...)
}

// Flush 由于底层 api 允许缓冲, 所以在进程退出之前调用 Sync
func Flush() {
	_ = logger.Sync()
}

func Debug(msg string) {
	logger.Debug(msg)
}

func Debugf(format string, v ...interface{}) {
	logger.Debug(fmt.Sprintf(format, v...))
}

func Info(msg string) {
	logger.Info(msg)
}

func Infof(format string, v ...interface{}) {
	logger.Info(fmt.Sprintf(format, v...))
}

func Warn(msg string) {
	logger.Warn(msg)
}

func Warnf(format string, v ...interface{}) {
	logger.Warn(fmt.Sprintf(format, v...))
}

func Error(msg string) {
	logger.Error(msg)
}

func Errorf(format string, v ...interface{}) {
	logger.Error(fmt.Sprintf(format, v...))
}

func Fatal(msg string) {
	logger.Fatal(msg)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatal(fmt.Sprintf(format, v...))
}
