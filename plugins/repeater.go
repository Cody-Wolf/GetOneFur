package plugins

import (
	"GetOneFur/messages"
	"GetOneFur/sender"
	"log"
	"strconv"
)

type Repeater struct {
	lastMessage string
	lastCount   int
}

func (r *Repeater) Response(message messages.Messager) {
	if r.lastMessage == message.GetMessage() {
		r.lastCount++
		log.Println("repeating message", r.lastMessage)
	} else {
		r.lastMessage = message.GetMessage()
		r.lastCount = 1
	}

	if r.lastCount == 3 {
		switch message.(type) {
		case messages.GroupMessager:
			groupMessager := message.(messages.GroupMessager)
			sender.SendGroupMessage(strconv.FormatInt(groupMessager.GetGroupId(), 10), groupMessager.GetMessage())
		default:
			log.Panicln("Not implemented message type. Failed repeated. Message = ", message)
		}
	}
}
