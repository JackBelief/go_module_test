package es_proc


type StudentInfo struct {
	Id         int `json:"id"`
	Title      string `json:"title"`
	Tags       string `json:"tags"`
	Short      string `json:"short"`
	Content    string `json:"content"`
	Author     string `json:"author"`
	CreateTime int64  `json:"create_time"`
}

type ESStudentInfo struct {
	EsIndex	string
	EsType	string
	EsData*	StudentInfo
}



