package messages

type GroupMessage struct {
	GroupId int64  `json:"group_id"`
	Message string `json:"message"`
}
