package gracefulshutdown

import (
	"os"
	"os/signal"
	"syscall"

	applogger "capstonea03/be/src/libs/logger"
)

type FnRunInShutdown struct {
	FnDescription string
	Fn            func()
}

var fnsRunInShutdown []FnRunInShutdown
var logger = applogger.New("GracefullShutdown")

func Add(newFns ...FnRunInShutdown) {
	fnsRunInShutdown = append(fnsRunInShutdown, newFns...)
}

func Run() {
	logger.Log("listening to shutdown signals")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGTERM)
	go func() {
		<-c
		if len(fnsRunInShutdown) > 0 {
			logger.Log("start clearing resources")
		}
		for i := range fnsRunInShutdown {
			logger.Log(fnsRunInShutdown[len(fnsRunInShutdown)-(i+1)].FnDescription)
			fnsRunInShutdown[len(fnsRunInShutdown)-(i+1)].Fn()
		}
		close(c)
	}()
}
