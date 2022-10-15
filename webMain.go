package main

import (
	"GetOneFur/messages"
	"GetOneFur/plugins"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var msgDB []string
var initPlugins []plugins.Plugin

func getBody(_ http.ResponseWriter, r *http.Request) {
	// 得到 request 的内容 json bytes
	size := r.ContentLength
	body := make([]byte, size)
	r.Body.Read(body)

	// 反序列化 json bytes 得到 map
	msgMap := make(map[string]interface{})
	if err := json.Unmarshal(body, &msgMap); err != nil {
		log.Panicln("反序列化错误, err = ", err)
		return
	}

	// TODO: 做一个 message parser
	var groupMessage messages.GroupMessage
	json.Unmarshal(body, &groupMessage)
	log.Println("group message : ", groupMessage)

	log.Printf("收到一条消息：group_id = %d \n消息内容：%s",
		groupMessage.GroupId, groupMessage.Message)

	for _, plugin := range initPlugins {
		plugin.Response(&groupMessage)
	}
}

func init() {
	initPlugins = append(initPlugins, new(plugins.RandomPicture))
	initPlugins = append(initPlugins, new(plugins.Repeater))
}

func main() {
	http.HandleFunc("/", getBody)
	// 在 5701 开一个 web 服务监听
	err := http.ListenAndServe(":5701", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
