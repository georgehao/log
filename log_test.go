package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDebug(t *testing.T) {
	Init("./test.log", DebugLevel, true)

	Debug("test test test")
	Debugf("test")
	Debugw("test", "a", "b")

	Info("test test test")
	Infof("test")
	Infow("test", "a", "b")

	Warn("test test test")
	Warnf("test")
	Warnw("test", "a", "b")

	Error("test test test")
	Errorf("test")
	Errorw("test", "a", "b")

	RequestLogInfow("aa", "bb")
	//Panic("test test test")
	//Panicf("test")
	//Panicw("test", "a", "b")
	//
	//Fatal("test")
	Sync()
}

func TestSetCaller(t *testing.T) {
	Init("./test.log", DebugLevel, true, SetCaller(true))
	assert.Equal(t, logger.isCaller(DebugLevel), true)
	Debug("test")
	Sync()
}
