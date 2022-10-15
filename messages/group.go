package messages

import (
	"encoding/json"
	"log"
)

type GroupMessage struct {
	GroupId int64  `json:"group_id"`
	Message string `json:"message"`
}

func (msg *GroupMessage) GetMessage() string {
	return msg.Message
}

func (msg *GroupMessage) GetGroupId() int64 {
	return msg.GroupId
}

func (msg *GroupMessage) SetMessage(msgBytes []byte) error {
	if err := json.Unmarshal(msgBytes, msg); err != nil {
		log.Panicln("反序列化错误, err = ", err)
		return err
	}
	return nil
}
