package main

import (
	"go_modules_test/src/nsq_proc"
)


func main() {
	if err := nsq_proc.InitNSQDClient(); err != nil {
		return
	}

	return
}
