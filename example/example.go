package main

import "github.com/georgehao/log"

func main() {
	// init log
	// set absolute path, and level
	// set output level
	// don't need request log
	// set log's caller using logOption
	log.Init("./test.log", log.DebugLevel, false, log.SetCaller(true))
	log.Info("hello george log")
	// flush
	log.Sync()
	//output: {"level":"info","ts":"2019-12-16T10:37:11.364+0800","caller":"example/example.go:12","msg":"hello george log"}
}