package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

// InitLogger 初始化 Zap 日志记录器，包含调用者信息
func InitLogger() {
	// 配置 JSON 输出的编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:    "timestamp",
		LevelKey:   "level",
		NameKey:    "logger",
		CallerKey:  "caller",
		MessageKey: "msg",
		// StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短格式：文件.go:行号
	}

	// 创建生产环境的配置
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// 使用 AddCaller 和 AddCallerSkip 构建日志记录器
	// AddCallerSkip(1) 跳过 InitLogger 函数，指向实际调用者
	logger, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	// 创建 SugaredLogger 以便于日志记录
	Logger = logger.Sugar()

	// 确保程序退出时刷新日志
	defer Logger.Sync()
}

// LogError 记录带有额外上下文的错误日志
func LogError(msg string, err error) {
	if err != nil {
		Logger.Errorw(msg,
			"error", err.Error(),
			// 如需更多上下文，可添加
			// "context", msg,
		)
	}
}
