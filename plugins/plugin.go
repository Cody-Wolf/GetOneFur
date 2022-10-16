package plugins

import "GetOneFur/messages"

type Plugin interface {
	Response(messages.Messager)
	HelpInfo() string
	GetPluginName() string
}
