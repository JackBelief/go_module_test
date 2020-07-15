package etcd_proc

import "go.etcd.io/etcd/clientv3"

var ETCDSchema = "ns"
var ServiceName = "say_hello_servers"

var ETCDServerPrefix = "/" + ETCDSchema + "/" + ServiceName + "/" // ETCD注册key前缀

var GClient *clientv3.Client // ETCD客户端
