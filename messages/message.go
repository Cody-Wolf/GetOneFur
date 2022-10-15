package messages

type Messager interface {
	GetMessage() string
	SetMessage([]byte) error
}

type GroupMessager interface {
	Messager
	GetGroupId() int64
}
