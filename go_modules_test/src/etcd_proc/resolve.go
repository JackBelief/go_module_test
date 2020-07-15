package etcd_proc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
)

/*****************************************************************************************************
	Builder 是接口类型，用于创建命名解析器，可监视命名空间是否发生变化，其方法有：
	1） Scheme() string		// 返回解析器支持的方案
	2） Build(target Target, cc ClientConn, opts BuildOptions) (Resolver, error)	// 创建解析器

	Resolver 是接口类型，用于监控目标变化，当目标发生变化时，会相应地更新地址、服务配置，其方法有：
	1） Close()		// 关闭解析器
	2） ResolveNow(ResolveNowOptions)		// 备用接口，GRPC可以再次调用用于目标的解析

	客户端要实现以上接口，从而实现服务发现、变更
*****************************************************************************************************/
func NewResolver() resolver.Builder {
	return &ETCDResolver{rawAddr: "121.42.161.154:2379"}
}

type ETCDResolver struct {
	rawAddr      string              // etcd服务地址，多个地址要使用分隔符
	resolverConn resolver.ClientConn // 解析器链接对象
}

// 实现Builder接口类型
func (er *ETCDResolver) Scheme() string {
	return ETCDSchema
}

func (er *ETCDResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// 构建解析器，解析器只负责对目标的更新，而对目标的监控由用户部分完成，
	var err error
	if GClient == nil {
		GClient, err = clientv3.New(clientv3.Config{
			Endpoints:   strings.Split(er.rawAddr, ";"),
			DialTimeout: 5 * time.Second,
		})

		if err != nil {
			return nil, err
		}
	}

	// 解析器监控变化
	er.resolverConn = cc
	fmt.Println("resolver create success")
	go er.watch("/" + target.Scheme + "/" + target.Endpoint + "/")

	return er, nil
}

func (er *ETCDResolver) watch(keyPrefix string) {
	for {
		er.watchETCD(keyPrefix)
		time.Sleep(1 * time.Second)
	}
}

func (er *ETCDResolver) watchETCD(keyPrefix string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("watch error =", err)
		}
	}()

	er.watchETCDKey(keyPrefix)
}

func (er *ETCDResolver) watchETCDKey(keyPrefix string) {
	var addrList []resolver.Address

	// 读取ETCD，获取IP列表
	getResp, err := GClient.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("解析器读取ETCD，获取IP列表失败 err=", err.Error())
	} else {
		for index := range getResp.Kvs {
			fmt.Println("初始IP地址是：", strings.TrimPrefix(string(getResp.Kvs[index].Key), keyPrefix))
			addrList = append(addrList, resolver.Address{Addr: strings.TrimPrefix(string(getResp.Kvs[index].Key), keyPrefix)})
		}
	}

	er.resolverConn.NewAddress(addrList)
	// er.resolverConn.UpdateState(resolver.State{Addresses:addrList})

	// 监控ETCD中目标数据的变化
	watchChan := GClient.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
	for chanEle := range watchChan {
		for _, ev := range chanEle.Events {
			// 根据IP变化情况，解析器更新IP地址列表
			addr := strings.TrimPrefix(string(ev.Kv.Key), keyPrefix)
			switch ev.Type {
			case mvccpb.PUT:
				if !exist(addrList, addr) {
					addrList = append(addrList, resolver.Address{Addr: addr})
					er.resolverConn.NewAddress(addrList)
					fmt.Println("插入新地址 address=", addr)
				}
			case mvccpb.DELETE:
				if s, ok := remove(addrList, addr); ok {
					addrList = s
					er.resolverConn.NewAddress(addrList)
					fmt.Println("删除老地址 address=", addr)
				}
			}
		}
	}
}

func exist(l []resolver.Address, addr string) bool {
	for i := range l {
		if l[i].Addr == addr {
			return true
		}
	}

	return false
}

func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}

	return nil, false
}

// 实现Resolver接口类型
func (er *ETCDResolver) ResolveNow(rn resolver.ResolveNowOptions) {
	fmt.Println("ETCDResolver ResolveNow")
}

func (er *ETCDResolver) Close() {
	fmt.Println("ETCDResolver Close")
}
