package main

import (
	"fmt"
	"go_modules_test/src/config"
	"go_modules_test/src/etcd_proc"
	"go_modules_test/src/protoes"
	"go_modules_test/src/service"
	"os"
	"os/signal"
	"syscall"

	"net"

	"google.golang.org/grpc"
)

func main() {
	// 解析命令行参数，获取服务启动地址和端口
	config.ParseServerCommandLine()

	// 向ETCD注册服务
	etcd_proc.RegisterETCDServer(config.GCfg.GetServerAddr())

	// 启动服务
	startServer()
}

func startServer() {
	// step 1 : 创建监听套接字
	serverAddr := config.GCfg.GetServerAddr()
	listen, err := net.Listen("tcp", serverAddr)
	if err != nil {
		fmt.Println("服务监听套接字创建失败 ip=", serverAddr, " error=", err.Error())
		return
	}

	// step 2 : 创建GRPC服务对象
	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	// step 3 : 向GRPC服务注册服务结构体
	protoes.RegisterHelloGrpcServer(grpcServer, &service.GRPCServiceTest{})

	// step 4 : 信号监听处理
	go signalProc()

	// step 5 : 启动服务
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Println("服务启动失败", err)
		return
	}
}

// 信号捕捉处理函数
func signalProc() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		chanData := <-signalChan
		etcd_proc.UnRegisterETCDServer(config.GCfg.GetServerAddr())

		if signalId, ok := chanData.(syscall.Signal); ok == true {
			os.Exit(int(signalId))
		} else {
			os.Exit(0)
		}
	}()
}
