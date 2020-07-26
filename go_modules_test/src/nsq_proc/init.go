package nsq_proc

import (
	//"encoding/json"
	//"fmt"

	"fmt"
	"github.com/nsqio/go-nsq"
)

var GNSQDClient *nsq.Producer

// 初始化生产者
func Init() error {
	producer, err := nsq.NewProducer(NSQDAddr, nsq.NewConfig())
	if err != nil {
		fmt.Println("nsqd client init fail ", err)
		return err
	}

	GNSQDClient = producer
	return nil
}
