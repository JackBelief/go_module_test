package service

import (
	"context"
	"fmt"
	"go_modules_test/src/config"
	"go_modules_test/src/protoes"
)

// GRPC服务结构要和proto内定义的服务结构保持一致
type GRPCServiceTest struct {
}

func (s *GRPCServiceTest) SayHi(ctx context.Context, in *protoes.HelloRequest) (*protoes.HelloReplay, error) {
	fmt.Println(config.GCfg.GetServerAddr(), " SayHi ", in.Name)
	return &protoes.HelloReplay{
		Message: "Hi, My name is " + in.Name,
	}, nil
}

func (s *GRPCServiceTest) GetMsg(ctx context.Context, in *protoes.HelloRequest) (*protoes.HelloMessage, error) {
	fmt.Println(config.GCfg.GetServerAddr(), " GetMsg ", in.Name)
	return &protoes.HelloMessage{
		Msg: "Msg is that " + in.Name + " is coming",
	}, nil
}
