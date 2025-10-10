package logger

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

var Logger2 *zap.SugaredLogger

// yamlEncoder 实现 zapcore.Encoder 接口以支持 YAML 格式
type yamlEncoder struct {
	config zapcore.EncoderConfig
}

// bufferEncoder 实现 zapcore.PrimitiveArrayEncoder 接口以捕获编码输出
type bufferEncoder struct {
	buf *bytes.Buffer
}

func (enc *bufferEncoder) AppendBool(v bool)              { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendByteString(v []byte)      { enc.buf.Write(v) }
func (enc *bufferEncoder) AppendComplex128(v complex128)  { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendComplex64(v complex64)    { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendFloat64(v float64)        { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendFloat32(v float32)        { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendInt(v int)                { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendInt64(v int64)            { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendInt32(v int32)            { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendInt16(v int16)            { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendInt8(v int8)              { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendString(v string)          { enc.buf.WriteString(v) }
func (enc *bufferEncoder) AppendUint(v uint)              { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendUint64(v uint64)          { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendUint32(v uint32)          { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendUint16(v uint16)          { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendUint8(v uint8)            { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendUintptr(v uintptr)        { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendDuration(v time.Duration) { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendTime(v time.Time)         { fmt.Fprintf(enc.buf, "%v", v) }
func (enc *bufferEncoder) AppendArray(v zapcore.ArrayMarshaler) error {
	return fmt.Errorf("嵌套数组未实现")
}
func (enc *bufferEncoder) AppendObject(v zapcore.ObjectMarshaler) error {
	return fmt.Errorf("嵌套对象未实现")
}
func (enc *bufferEncoder) AppendReflected(v interface{}) error {
	fmt.Fprintf(enc.buf, "%v", v)
	return nil
}

// EncodeEntry 将日志条目编码为 YAML 格式
func (enc *yamlEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	logEntry := make(map[string]interface{})

	// 编码时间
	if enc.config.TimeKey != "" && enc.config.EncodeTime != nil {
		var buf bytes.Buffer
		enc.config.EncodeTime(entry.Time, &bufferEncoder{buf: &buf})
		logEntry[enc.config.TimeKey] = buf.String()
	}

	// 编码日志级别
	if enc.config.LevelKey != "" && enc.config.EncodeLevel != nil {
		var buf bytes.Buffer
		enc.config.EncodeLevel(entry.Level, &bufferEncoder{buf: &buf})
		logEntry[enc.config.LevelKey] = buf.String()
	}

	// 编码调用者信息
	if enc.config.CallerKey != "" && enc.config.EncodeCaller != nil && entry.Caller.Defined {
		var buf bytes.Buffer
		enc.config.EncodeCaller(entry.Caller, &bufferEncoder{buf: &buf})
		logEntry[enc.config.CallerKey] = buf.String()
	}

	// 编码消息
	if enc.config.MessageKey != "" {
		logEntry[enc.config.MessageKey] = entry.Message
	}

	// 添加额外的字段
	for _, field := range fields {
		// 调试：打印字段值以检查是否为 nil
		if field.Key == "error" && field.Interface == nil {
			logEntry[field.Key] = "nil-error"
		} else if field.Interface == "" {
			logEntry[field.Key] = "empty-string"
		} else {
			logEntry[field.Key] = field.Interface
		}
	}

	// 转换为 YAML 格式
	buf, err := yaml.Marshal(logEntry)
	if err != nil {
		return nil, fmt.Errorf("无法序列化 YAML: %v", err)
	}

	// 创建缓冲区并写入 YAML 数据
	buffer := buffer.NewPool().Get()
	if _, err := buffer.Write(buf); err != nil {
		return nil, fmt.Errorf("无法写入缓冲区: %v", err)
	}
	buffer.AppendString("---\n") // YAML 文档分隔符
	return buffer, nil
}

// Clone 实现 Encoder 接口的 Clone 方法
func (enc *yamlEncoder) Clone() zapcore.Encoder {
	return &yamlEncoder{config: enc.config}
}

// 实现 zapcore.Encoder 接口的所有方法
func (enc *yamlEncoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error   { return nil }
func (enc *yamlEncoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error { return nil }
func (enc *yamlEncoder) AddBinary(key string, val []byte)                              {}
func (enc *yamlEncoder) AddByteString(key string, val []byte)                          {}
func (enc *yamlEncoder) AddBool(key string, val bool)                                  {}
func (enc *yamlEncoder) AddComplex128(key string, val complex128)                      {}
func (enc *yamlEncoder) AddComplex64(key string, val complex64)                        {}
func (enc *yamlEncoder) AddDuration(key string, val time.Duration)                     {}
func (enc *yamlEncoder) AddFloat64(key string, val float64)                            {}
func (enc *yamlEncoder) AddFloat32(key string, val float32)                            {}
func (enc *yamlEncoder) AddInt(key string, val int)                                    {}
func (enc *yamlEncoder) AddInt64(key string, val int64)                                {}
func (enc *yamlEncoder) AddInt32(key string, val int32)                                {}
func (enc *yamlEncoder) AddInt16(key string, val int16)                                {}
func (enc *yamlEncoder) AddInt8(key string, val int8)                                  {}
func (enc *yamlEncoder) AddString(key string, val string)                              {}
func (enc *yamlEncoder) AddTime(key string, val time.Time)                             {}
func (enc *yamlEncoder) AddUint(key string, val uint)                                  {}
func (enc *yamlEncoder) AddUint64(key string, val uint64)                              {}
func (enc *yamlEncoder) AddUint32(key string, val uint32)                              {}
func (enc *yamlEncoder) AddUint16(key string, val uint16)                              {}
func (enc *yamlEncoder) AddUint8(key string, val uint8)                                {}
func (enc *yamlEncoder) AddUintptr(key string, val uintptr)                            {}
func (enc *yamlEncoder) AddReflected(key string, val interface{}) error                { return nil }
func (enc *yamlEncoder) OpenNamespace(key string)                                      {}

// NewYAMLEncoder 创建一个新的 YAML 编码器
func NewYAMLEncoder(config zapcore.EncoderConfig) zapcore.Encoder {
	return &yamlEncoder{config: config}
}

// InitLogger 初始化 Zap 日志记录器，输出到 ./yaml/时间戳.yaml 文件并支持日志轮转
func InitLogger2() {
	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短格式：文件.go:行号
	}

	// 确保 ./yaml 目录存在
	logDir := "./yaml"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Errorf("无法创建日志目录 %s: %v", logDir, err))
	}

	// 生成时间戳格式的日志文件名
	logFileName := time.Now().Format("2006-01-02/15-04-05") + ".yaml"
	logFilePath := filepath.Join(logDir, logFileName)

	// 配置日志轮转
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFilePath, // 初始日志文件
		MaxSize:    100,         // 最大 100 MB
		MaxAge:     28,          // 最大保留 28 天
		MaxBackups: 0,           // 不限制备份数量（仅受 MaxAge 限制）
		Compress:   false,       // 不压缩备份文件
	})

	// 配置 Zap 核心
	core := zapcore.NewCore(
		NewYAMLEncoder(encoderConfig), // 使用自定义 YAML 编码器
		writer,                        // 使用 lumberjack 进行日志轮转
		zapcore.DebugLevel,            // 日志级别
	)

	// 构建日志记录器，包含调用者信息
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Logger = logger.Sugar()

	// 确保程序退出时刷新日志
	defer Logger.Sync()
}

// LogError 记录带有额外上下文的错误日志
func LogError2(msg string, err error) {
	if err != nil {
		Logger.Errorw(msg,
			"error", err.Error(), // 确保 err.Error() 的值被正确传递
			// "context", msg, // 可选上下文，当前注释掉
		)
	} else {
		// 如果 err 为 nil，记录一条调试日志
		Logger.Debugw("错误为空，未记录错误日志",
			"msg", msg,
		)
	}
}
