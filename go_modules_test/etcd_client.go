package main

import (
	"context"
	"fmt"
	"go_modules_test/src/config"
	"go_modules_test/src/etcd_proc"
	"go_modules_test/src/protoes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"os"
	"os/signal"
	"time"
)

func main() {
	defer func() {
		c := make(chan os.Signal)
		signal.Notify(c)
		s := <-c
		fmt.Println("get signal:", s)
	}()

	// 解析命令行参数，获取服务启动地址和端口
	config.ParseServerCommandLine()

	// 创建命名解析
	r := etcd_proc.NewResolver()
	resolver.Register(r)

	conn, err := grpc.Dial(r.Scheme()+"://author/"+etcd_proc.ServiceName, grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// 获得grpc句柄，定时发送消息
	c := protoes.NewHelloGrpcClient(conn)
	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		// 远程单调用 SayHi 接口
		r1, err := c.SayHi(
			context.Background(),
			&protoes.HelloRequest{
				Name: "Kitty",
			},
		)

		if err != nil {
			fmt.Println("命名解析异常，没有服务地址 err=", err)
			continue
		}

		fmt.Printf("%v: SayHi 响应：%s\n", t, r1.GetMessage())

		// 远程单调用 GetMsg 接口
		r2, err := c.GetMsg(
			context.Background(),
			&protoes.HelloRequest{
				Name: "Jack",
			},
		)

		if err != nil {
			fmt.Println("命名解析异常，没有服务地址 err=", err)
			continue
		}

		fmt.Printf("%v: GetMsg 响应：%s\n", t, r2.GetMsg())
	}

	fmt.Println("client over")
	return
}
