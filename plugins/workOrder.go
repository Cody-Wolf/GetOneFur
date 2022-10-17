package plugins

import "GetOneFur/messages"

type workOrder struct{}

func (r *workOrder) Response(message messages.Messager) {

	//TODO implement me
	panic("implement me")
}

func (r *workOrder) HelpInfo() string {
	return "/提问：选择需要提问的消息并回复，回复 “/提问” 即可提问。\n" +
		"/回答：\n"
}

func (r *workOrder) GetPluginName() string {
	return "工单处理"
}
