package etcd_proc

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func RegisterETCDServer(addr string) {
	// 服务注册
	registerServer(addr)
}

func registerServer(addr string) {
	var err error

	// 创建ETCD的客户端
	if GClient == nil {
		GClient, err = newETCDClient()
		if err != nil {
			fmt.Println("ectd 客户端创建失败 error=", err.Error())
			return
		}
	}
	fmt.Println("ectd 客户端创建成功")

	// 定时循环检测，查看向ETCD注册服务是否正常
	// 每台服务向ETCD注册自己的IP地址，定时检测注册内容是否还在
	ticker := time.NewTicker(time.Second * time.Duration(5))
	go func() {
		for {
			getResp, err := GClient.Get(context.Background(), ETCDServerPrefix+addr)
			if err != nil {
				fmt.Println("etcd出现异常，key获取异常，key=", ETCDServerPrefix+addr, " error=", err.Error())
			} else if getResp.Count == 0 {
				fmt.Println("etcd没有目标数据，需要补数据，key=", ETCDServerPrefix+addr)
				putData(ETCDServerPrefix+addr, addr)
			} else {
				fmt.Println("etcd目标数据正常，key=", ETCDServerPrefix+addr)
			}

			<-ticker.C
		}
	}()

	return
}

func newETCDClient() (*clientv3.Client, error) {
	config := clientv3.Config{
		Endpoints:   []string{"121.42.161.154:2379"},
		DialTimeout: 5 * time.Second,
	}

	return clientv3.New(config)
}

func putData(key, value string) {
	leaseResp, err := GClient.Grant(context.Background(), 5)
	if err != nil {
		fmt.Println("etcd申请租约失败 key=", key, " error=", err.Error())
		return
	}

	_, err = GClient.Put(context.Background(), key, value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		fmt.Println("etcd写入数据失败 key=", key, " error=", err.Error())
		return
	}

	kaRespChan, err := GClient.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		fmt.Println("etcd租约续约失败 key=", key, "id=", leaseResp.ID, " error=", err.Error())
		return
	}

	// 定期查看续约结果
	go func() {
		for {
			select {
			case respData := <-kaRespChan:
				if kaRespChan == nil {
					fmt.Println("管道关闭，出现异常，退出 key=", key)
					return
				} else {
					if respData == nil {
						fmt.Println("没有数据，可能是etcd关闭、也可能是网络异常，退出 key=", key)
						return
					} else {
						fmt.Println("续约成功 key=", key)
					}
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return
}

func UnRegisterETCDServer(addr string) {
	// 服务取消注册
	unRegisterServer(addr)
}

func unRegisterServer(addr string) {
	var err error

	// 创建ETCD的客户端
	if GClient == nil {
		GClient, err = newETCDClient()
		if err != nil {
			fmt.Println("ectd 客户端创建失败 error=", err.Error())
			return
		}
	}

	// 删除服务注册数据
	_, err = GClient.Delete(context.Background(), ETCDServerPrefix+addr)
	if err != nil {
		fmt.Println("服务关闭，etcd删除数据失败 key=", ETCDServerPrefix+addr, " error=", err.Error())
		return
	} else {
		fmt.Println("服务关闭，etcd成功删除数据 key=", ETCDServerPrefix+addr)
		return
	}
}
