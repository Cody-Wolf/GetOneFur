package plugins

import (
	"GetOneFur/messages"
	"GetOneFur/sender"
	"log"
	"strconv"
)

type StringCount struct {
	str   string
	count int
}

type Repeater struct {
	repeatMap map[int64]StringCount
}

func (r *Repeater) HelpInfo() string {
	return "每个群达到 3 次的重复消息就会被复读一次。再次重复不会复读。\n"
}

func (r *Repeater) GetPluginName() string {
	return "复读机"
}

func (r *Repeater) Response(message messages.Messager) {
	if r.repeatMap == nil {
		r.repeatMap = make(map[int64]StringCount)
	}

	switch message.(type) {
	case messages.GroupMessager:
		groupMessager := message.(messages.GroupMessager)
		if groupMessager.GetGroupId() == 0 {
			break
		}
		log.Println("message = ", message)
		stringCount, ok := r.repeatMap[groupMessager.GetGroupId()]
		if ok && stringCount.str == message.GetMessage() {
			stringCount.count++
		} else {
			stringCount = StringCount{groupMessager.GetMessage(), 1}
		}
		if stringCount.count == 3 {
			sender.SendGroupMessage(strconv.FormatInt(groupMessager.GetGroupId(), 10), groupMessager.GetMessage())
		}
		r.repeatMap[groupMessager.GetGroupId()] = stringCount
	default:
		log.Panicln("不能识别的消息类型，复读失败。")
	}
}
