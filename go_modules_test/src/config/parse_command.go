package config

import (
	"flag"
	"strconv"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// 服务启动参数配置，用于单机上模拟多服务器的情况
var GCfg *ServerConfig

type ServerConfig struct {
	Ip   string
	Port int
}

func (sc *ServerConfig) GetServerAddr() string {
	return sc.Ip + ":" + strconv.Itoa(sc.Port)
}

func ParseServerCommandLine() {
	flag.Int("port", 1688, "server port")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	cfg := &ServerConfig{
		Ip: "127.0.0.1",
	}
	cfg.Port = viper.GetInt("port")
	GCfg = cfg

	return
}
