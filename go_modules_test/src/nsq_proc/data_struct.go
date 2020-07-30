package nsq_proc

// 发布结构
type StudentPubMsg struct {
	Id int `json:"id"`
	Action int `json:"action"`
}
