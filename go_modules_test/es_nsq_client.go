package main

import (
	"fmt"
	"go_modules_test/src/es_proc"
	"go_modules_test/src/nsq_proc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if InitClient() != nil {
		return
	}

	EsNsqTestClient()
	return
}

func InitClient()(err error) {
	// es初始化
	if err = es_proc.Init(); err != nil {
		fmt.Println("es init fail", err)
		return err
	}

	// nsq初始化
	if err = nsq_proc.InitNSQDClient(); err != nil {
		fmt.Println("nsq init fail", err)
		return err
	}

	return err
}

func EsNsqTestClient() {
	var err error
	if err = nsq_proc.InitNSQDClient(); err != nil {
		fmt.Println(err)
		return
	}

	if err = nsq_proc.PutPubTopic(nsq_proc.StudentPubTopic, 1); err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(1*time.Second)
	if err = nsq_proc.DeletePubTopic(nsq_proc.StudentPubTopic, 1); err != nil {
		fmt.Println(err)
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	<-ch
}