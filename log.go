package log

import (
	"fmt"
)

const (
	DebugLevel   = "debug"
	InfoLevel    = "info"
	WarnLevel    = "warn"
	ErrorLevel   = "error"
	PanicLevel   = "panic"
	RequestLevel = "request"
)

var logger *Log

// Log 默认会使用zap作为日志输出引擎. Log集成了日志切割的功能。默认文件大小1024M，自动压缩
// 最大有3个文件备份，备份保存时间7天。默认不会打印日志被调用的文文件名和位置;
// 输出:日志默认会被分成五类文件：xxx.log.DEBUG，xxx.log.INFO, xxx.log.WARN, xxx.log.ERROR, xxx.log.Request
// error,panic都会打印在xxx.log.ERROR. 所有的请求都会打在xxx.log.Request
// Adapter:经过比对现在流行的日志库：zap, logrus, zerolog; logrus 虽说格式化，插件化良好，但是
// 其内部实现锁竞争太过剧烈，性能不好. zap 性能好，格式一般， zerolog性能没有zap好， 相比
// 来说就没啥突出优点了

type Log struct {
	Path           string
	Level          string
	NeedRequestLog bool // 是否需要request.log
	adapters       map[string]*zapAdapter
}

type LogOption interface {
	apply(*Log)
}

type logOptionFunc func(*Log)

func (f logOptionFunc) apply(log *Log) {
	f(log)
}

func SetMaxFileSize(size int) LogOption {
	return logOptionFunc(func(log *Log) {
		for i, _ := range log.adapters {
			log.adapters[i].setMaxFileSize(size)
		}
	})
}

func SetMaxBackups(n int) LogOption {
	return logOptionFunc(func(log *Log) {
		for i, _ := range log.adapters {
			log.adapters[i].setMaxBackups(n)
		}
	})
}

func SetMaxAge(age int) LogOption {
	return logOptionFunc(func(log *Log) {
		for i, _ := range log.adapters {
			log.adapters[i].setMaxAge(age)
		}
	})
}

func SetCompress(compress bool) LogOption {
	return logOptionFunc(func(log *Log) {
		for i, _ := range log.adapters {
			log.adapters[i].setCompress(compress)
		}
	})
}

func SetCaller(caller bool) LogOption {
	return logOptionFunc(func(log *Log) {
		for i, _ := range log.adapters {
			log.adapters[i].setCaller(caller)
		}
	})
}

// Init init logger
func Init(path, level string, needRequestLog bool, options ...LogOption) {
	logger = &Log{Path: path, Level: level}

	var levels []string
	switch level {
	case DebugLevel:
		levels = append(levels, DebugLevel, InfoLevel, WarnLevel, ErrorLevel)
	case InfoLevel:
		levels = append(levels, InfoLevel, WarnLevel, ErrorLevel)
	case WarnLevel:
		levels = append(levels, WarnLevel, ErrorLevel)
	case ErrorLevel, PanicLevel:
		levels = append(levels, ErrorLevel)
	}

	if needRequestLog {
		logger.NeedRequestLog = true
		levels = append(levels, RequestLevel)
	}

	logger.createFiles(levels, options...)
}

// Sync flushes buffer, if any
func Sync() {
	if logger == nil {
		return
	}

	for _, v := range logger.adapters {
		v.logger.Sync()
	}
}

func (l *Log) isCaller(level string) bool {
	if v, ok := l.adapters[level]; ok {
		return v.Caller
	}
	return false
}

func (l *Log) maxFileSize(level string) int {
	if v, ok := l.adapters[level]; ok {
		return v.MaxFileSize
	}
	return 0
}

func (l *Log) maxBackups(level string) int {
	if v, ok := l.adapters[level]; ok {
		return v.MaxBackups
	}
	return 0
}

func (l *Log) maxAge(level string) int {
	if v, ok := l.adapters[level]; ok {
		return v.MaxAge
	}
	return 0
}

func (l *Log) createFiles(levels []string, options ...LogOption) {
	adapters := make(map[string]*zapAdapter, 4)
	for _, v := range levels {
		switch {
		case v == DebugLevel:
			adapters[DebugLevel] = NewZapAdapter(fmt.Sprintf("%s.DEBUG", l.Path), DebugLevel)
		case v == InfoLevel:
			adapters[InfoLevel] = NewZapAdapter(fmt.Sprintf("%s.INFO", l.Path), InfoLevel)
		case v == WarnLevel:
			adapters[WarnLevel] = NewZapAdapter(fmt.Sprintf("%s.WARN", l.Path), WarnLevel)
		case v == ErrorLevel:
			adapters[ErrorLevel] = NewZapAdapter(fmt.Sprintf("%s.ERROR", l.Path), ErrorLevel)
		case v == RequestLevel:
			adapters[RequestLevel] = NewZapAdapter(fmt.Sprintf("%s.Request", l.Path), InfoLevel)
		}
	}
	l.adapters = adapters

	for _, opt := range options {
		opt.apply(l)
	}

	for _, adapter := range adapters {
		adapter.Build()
	}

}

// Debug 使用方法：log.Debug("test")
func Debug(args ...interface{}) {
	if logger == nil {
		return
	}

	if debugAdapter, ok := logger.adapters[DebugLevel]; ok {
		debugAdapter.Debug(args...)
	}
}

// Debugf 使用方法：log.Debugf("test:%s", err)
func Debugf(template string, args ...interface{}) {
	if logger == nil {
		return
	}

	if debugAdapter, ok := logger.adapters[DebugLevel]; ok {
		debugAdapter.Debugf(template, args...)
	}
}

// Debugw 使用方法：log.Debugw("test", "field1", "value1", "field2", "value2")
func Debugw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	if debugAdapter, ok := logger.adapters[DebugLevel]; ok {
		debugAdapter.Debugw(msg, keysAndValues...)
	}
}

func Info(args ...interface{}) {
	if logger == nil {
		return
	}

	if adapter, ok := logger.adapters[InfoLevel]; ok {
		adapter.Info(args...)
	}
}

func Infof(template string, args ...interface{}) {
	if logger == nil {
		return
	}

	if adapter, ok := logger.adapters[InfoLevel]; ok {
		adapter.Infof(template, args...)
	}
}

func Infow(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	if adapter, ok := logger.adapters[InfoLevel]; ok {
		adapter.Infow(msg, keysAndValues...)
	}
}

func Warn(args ...interface{}) {
	if logger == nil {
		return
	}

	if adapter, ok := logger.adapters[WarnLevel]; ok {
		adapter.Warn(args...)
	}
}

func Warnf(template string, args ...interface{}) {
	if logger == nil {
		return
	}

	if adapter, ok := logger.adapters[WarnLevel]; ok {
		adapter.Warnf(template, args...)
	}
}

func Warnw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	if adapter, ok := logger.adapters[WarnLevel]; ok {
		adapter.Warnw(msg, keysAndValues...)
	}
}

func Error(args ...interface{}) {
	if logger == nil {
		return
	}

	if adapter, ok := logger.adapters[ErrorLevel]; ok {
		adapter.Error(args...)
	}
}

func Errorf(template string, args ...interface{}) {
	if logger == nil {
		return
	}

	if adapter, ok := logger.adapters[ErrorLevel]; ok {
		adapter.Errorf(template, args...)
	}
}

func Errorw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	if adapter, ok := logger.adapters[ErrorLevel]; ok {
		adapter.Errorw(msg, keysAndValues...)
	}
}

func Panic(args ...interface{}) {
	if logger == nil {
		return
	}

	if adapter, ok := logger.adapters[ErrorLevel]; ok {
		adapter.Panic(args...)
	}
}

func Panicf(template string, args ...interface{}) {
	if logger == nil {
		return
	}

	if adapter, ok := logger.adapters[ErrorLevel]; ok {
		adapter.Panicf(template, args...)
	}
}

func Panicw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	if adapter, ok := logger.adapters[ErrorLevel]; ok {
		adapter.Panicw(msg, keysAndValues...)
	}
}

func Fatal(args ...interface{}) {
	if logger == nil {
		return
	}

	if adapter, ok := logger.adapters[ErrorLevel]; ok {
		adapter.Fatal(args...)
	}
}

func Fatalf(template string, args ...interface{}) {
	if logger == nil {
		return
	}

	if adapter, ok := logger.adapters[ErrorLevel]; ok {
		adapter.Fatalf(template, args...)
	}
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	if adapter, ok := logger.adapters[ErrorLevel]; ok {
		adapter.Fatalw(msg, keysAndValues...)
	}
}

func RequestLogInfow(keysAndValues ...interface{}) {
	if logger == nil || !logger.NeedRequestLog {
		return
	}

	if adapter, ok := logger.adapters[RequestLevel]; ok {
		adapter.Infow("", keysAndValues...)
	}
}
