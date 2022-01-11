package logx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
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
// 	创建完成后不要忘记 defer xxx.Sync()
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

func Debug(msg string) {
	defer logger.Sync()
	logger.Debug(msg)
}

func Info(msg string) {
	defer logger.Sync()
	logger.Info(msg)

}

func Warn(msg string) {
	defer logger.Sync()
	logger.Warn(msg)
}

func Error(msg string) {
	defer logger.Sync()
	logger.Error(msg)
}

func Fatal(msg string) {
	defer logger.Sync()
	logger.Fatal(msg)
}
