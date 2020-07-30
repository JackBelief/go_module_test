package nsq_proc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"go_modules_test/src/es_proc"
)

type NSQDMetricsConsumer struct {

}

func (consumer* NSQDMetricsConsumer)HandleMessage(msg *nsq.Message) error {
	fmt.Println("consumer收到的消费信息是 addr:", msg.NSQDAddress, "message:", string(msg.Body))

	var data StudentPubMsg
	err := json.Unmarshal([]byte(msg.Body), &data)
	if err != nil {
		fmt.Println("nsq message element parse fail ", err)
		return err
	}

	switch data.Action {
	case StudentPublishActionPut:
		PutStudentInfo(data.Id)
	//case StudentPublishActionGet:
	//	GetStudentInfo(data.Id)
	case StudentPublishActionDelete:
		DeleteStudentInfo(data.Id)
	default:
		fmt.Println("unknow queue action: %d", data.Action)
	}

	return nil
}

var StudentInfoArray = []es_proc.ESStudentInfo{
	{"es_stu_info", "language", &es_proc.StudentInfo{1, "Go语言学习", "Go语言", "Go是一门很好的语言","Go语言使用场景很多", "田庆阳", 100}},
	{"es_stu_info", "language", &es_proc.StudentInfo{2, "C++语言学习", "C++语言", "C++是一门很好的语言","C++语言使用场景很多", "田庆钛", 200}},
	{"es_stu_info", "language", &es_proc.StudentInfo{3, "Python语言学习", "Python语言", "Python是一门很好的语言","Python语言使用场景很多", "田海洋", 300}},
}

func PutStudentInfo(id int) error {
	for _, ele := range StudentInfoArray {
		if ele.EsData.Id == id {
			return es_proc.PutData(context.Background(), &ele)
		}
	}

	fmt.Println("put no find student id=", id)
	return nil
}

func DeleteStudentInfo(id int) error {
	for _, ele := range StudentInfoArray {
		if ele.EsData.Id == id {
			return es_proc.DeleteData(context.Background(), &ele)
		}
	}

	fmt.Println("delete no find student id=", id)
	return nil
}