package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapAdapter struct {
	Path        string // 文件绝对地址，如：/home/homework/neso/file.log
	Level       string // 日志输出的级别
	MaxFileSize int    // 日志文件大小的最大值，单位(M)
	MaxBackups  int    // 最多保留备份数
	MaxAge      int    // 日志文件保存的时间，单位(天)
	Compress    bool   // 是否压缩
	Caller      bool   // 日志是否需要显示调用位置

	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func (z *zapAdapter) setMaxFileSize(size int) {
	z.MaxFileSize = size
}

func (z *zapAdapter) setMaxBackups(n int) {
	z.MaxBackups = n
}

func (z *zapAdapter) setMaxAge(age int) {
	z.MaxAge = age
}

func (z *zapAdapter) setCompress(compress bool) {
	z.Compress = compress
}

func (z *zapAdapter) setCaller(caller bool) {
	z.Caller = caller
}

func NewZapAdapter(path, level string) *zapAdapter {
	return &zapAdapter{
		Path:        path,
		Level:       level,
		MaxFileSize: 1024,
		MaxBackups:  3,
		MaxAge:      7,
		Compress:    true,
		Caller:      false,
	}
}

// createLumberjackHook 创建LumberjackHook，其作用是为了将日志文件切割，压缩
func (zapAdapter *zapAdapter) createLumberjackHook() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   zapAdapter.Path,
		MaxSize:    zapAdapter.MaxFileSize,
		MaxBackups: zapAdapter.MaxBackups,
		MaxAge:     zapAdapter.MaxAge,
		Compress:   zapAdapter.Compress,
	}
}

func (zapAdapter *zapAdapter) Build() {
	w := zapcore.AddSync(zapAdapter.createLumberjackHook())

	var level zapcore.Level
	switch zapAdapter.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "panic":
		level = zap.PanicLevel
	default:
		level = zap.InfoLevel
	}

	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	cnf := zapcore.NewJSONEncoder(conf)
	core := zapcore.NewCore(cnf, w, level)

	zapAdapter.logger = zap.New(core)
	if zapAdapter.Caller {
		zapAdapter.logger = zapAdapter.logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(2))
	}
	zapAdapter.sugar = zapAdapter.logger.Sugar()
}

func (zapAdapter *zapAdapter) Debug(args ...interface{}) {
	zapAdapter.sugar.Debug(args...)
}

func (zapAdapter *zapAdapter) Info(args ...interface{}) {
	zapAdapter.sugar.Info(args...)
}

func (zapAdapter *zapAdapter) Warn(args ...interface{}) {
	zapAdapter.sugar.Warn(args...)
}

func (zapAdapter *zapAdapter) Error(args ...interface{}) {
	zapAdapter.sugar.Error(args...)
}

func (zapAdapter *zapAdapter) DPanic(args ...interface{}) {
	zapAdapter.sugar.DPanic(args...)
}

func (zapAdapter *zapAdapter) Panic(args ...interface{}) {
	zapAdapter.sugar.Panic(args...)
}

func (zapAdapter *zapAdapter) Fatal(args ...interface{}) {
	zapAdapter.sugar.Fatal(args...)
}

func (zapAdapter *zapAdapter) Debugf(template string, args ...interface{}) {
	zapAdapter.sugar.Debugf(template, args...)
}

func (zapAdapter *zapAdapter) Infof(template string, args ...interface{}) {
	zapAdapter.sugar.Infof(template, args...)
}

func (zapAdapter *zapAdapter) Warnf(template string, args ...interface{}) {
	zapAdapter.sugar.Warnf(template, args...)
}

func (zapAdapter *zapAdapter) Errorf(template string, args ...interface{}) {
	zapAdapter.sugar.Errorf(template, args...)
}

func (zapAdapter *zapAdapter) DPanicf(template string, args ...interface{}) {
	zapAdapter.sugar.DPanicf(template, args...)
}

func (zapAdapter *zapAdapter) Panicf(template string, args ...interface{}) {
	zapAdapter.sugar.Panicf(template, args...)
}

func (zapAdapter *zapAdapter) Fatalf(template string, args ...interface{}) {
	zapAdapter.sugar.Fatalf(template, args...)
}

func (zapAdapter *zapAdapter) Debugw(msg string, keysAndValues ...interface{}) {
	zapAdapter.sugar.Debugw(msg, keysAndValues...)
}

func (zapAdapter *zapAdapter) Infow(msg string, keysAndValues ...interface{}) {
	zapAdapter.sugar.Infow(msg, keysAndValues...)
}

func (zapAdapter *zapAdapter) Warnw(msg string, keysAndValues ...interface{}) {
	zapAdapter.sugar.Warnw(msg, keysAndValues...)
}

func (zapAdapter *zapAdapter) Errorw(msg string, keysAndValues ...interface{}) {
	zapAdapter.sugar.Errorw(msg, keysAndValues...)
}

func (zapAdapter *zapAdapter) DPanicw(msg string, keysAndValues ...interface{}) {
	zapAdapter.sugar.DPanicw(msg, keysAndValues...)
}

func (zapAdapter *zapAdapter) Panicw(msg string, keysAndValues ...interface{}) {
	zapAdapter.sugar.Panicw(msg, keysAndValues...)
}

func (zapAdapter *zapAdapter) Fatalw(msg string, keysAndValues ...interface{}) {
	zapAdapter.sugar.Fatalw(msg, keysAndValues...)
}
