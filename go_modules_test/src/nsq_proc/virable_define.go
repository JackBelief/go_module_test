package nsq_proc

var NSQDAddr = "127.0.0.1:4150"
var NSQLookupDAddr = "127.0.0.1:4161"

var StudentPubTopic = "stuInfo"


// 订阅动作枚举
const (
	StudentPublishActionPut int = 1
	StudentPublishActionGet int = 2
	StudentPublishActionDelete int = 3
)