package nsq_proc

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"time"
)

// 初始化生产者
var GNSQDClient *nsq.Producer
func InitNSQDClient() error {
	producer, err := nsq.NewProducer(NSQDAddr, nsq.NewConfig())
	if err != nil {
		fmt.Println("nsqd client init fail ", err)
		return err
	}

	GNSQDClient = producer
	return nil
}

// 初始化消费者
func InitNSQDConsumer() error {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = 1 * time.Second

	consumer, err := nsq.NewConsumer(StudentPubTopic, StudentPubTopic+"_"+"metrics", cfg)
	if err != nil {
		panic(err)
	}

	consumer.SetLogger(nil, 0)
	consumer.AddHandler(&NSQDMetricsConsumer{})

	if err := consumer.ConnectToNSQLookupd(NSQLookupDAddr); err != nil {
		fmt.Println("nsq lookup connect fail ", err)
		return err
	}

	return nil
}