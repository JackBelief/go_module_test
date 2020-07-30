package main

import (
	"go_modules_test/src/nsq_proc"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := nsq_proc.InitNSQDConsumer(); err != nil {
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	<-ch
}
