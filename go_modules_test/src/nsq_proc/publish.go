package nsq_proc

import (
	"encoding/json"
	"errors"
	"fmt"
)

func PutPubTopic(topic string, id int) error {
	return SendTopicMessage(topic, id, StudentPublishActionPut)
}

func GetPubTopic(topic string, id int) error {
	return SendTopicMessage(topic, id, StudentPublishActionGet)
}

func DeletePubTopic(topic string, id int) error {
	return SendTopicMessage(topic, id, StudentPublishActionDelete)
}

func SendTopicMessage(topic string, id int, action int) error {
	topicMsg := &StudentPubMsg{
		Id:id,
		Action:action,
	}

	topicJson, err := json.Marshal(topicMsg)
	if err != nil {
		fmt.Println("topic to json fail ", err)
		return err
	}
	fmt.Println("topic to json success ", err)

	return publishTopicMessage(topic, topicJson)
}

func publishTopicMessage(topic string, msg []byte) error {
	if len(msg) == 0 {
		return errors.New("input msg is empty")
	}

	return GNSQDClient.Publish(topic, msg)
}
