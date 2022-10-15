package main

import (
	"GetOneFur/messages"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var msgDB []string
var groupSet map[string]bool

func sentMessage(opt string, msgMap url.Values) []byte {
	resp, err := http.PostForm("http://127.0.0.1:5700/"+opt, msgMap)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	return body
}

func sentGroupMessage(groupID, message string) []byte {
	return sentMessage("send_group_msg", url.Values{"group_id": {groupID}, "message": {message}})
}

func groupManage(groupID string, msgMap map[string]interface{}) {
	if msgMap["message"] == "114514" {
		sentGroupMessage(groupID, "19191810")
	}
	if msgMap["message"] == "乐" {
		sentGroupMessage(groupID, "乐")
	}
	if strings.Contains(msgMap["message"].(string), "[CQ:image") {
		//sentGroupMessage(groupID, "检测到一张图片")
	}
	if strings.Contains(msgMap["message"].(string), "福瑞 查看") {
		sentGroupMessage(groupID, "福瑞控是吧")
	}
	if strings.Contains(msgMap["message"].(string), "福瑞浓度") {
		sentGroupMessage(groupID, "当前图库一共有"+strconv.Itoa(len(msgDB))+"只福瑞")
	}
	if strings.Contains(msgMap["message"].(string), "[CQ:reply") &&
		strings.Contains(msgMap["message"].(string), "添加福瑞") {
		//sentGroupMessage(groupID, "检测到回复")
		msgID := strings.TrimPrefix(msgMap["message"].(string), "[CQ:reply,id=")
		msgID, _, _ = strings.Cut(msgID, "]")
		fmt.Println("回复消息 ID : ", msgID)
		resBody := sentMessage("get_msg", url.Values{"message_id": {msgID}})
		resMap := make(map[string]interface{})
		if err := json.Unmarshal(resBody, &resMap); err != nil {
			fmt.Printf("反序列化错误 err=%v\n", err)
			return
		}
		dataMap := resMap["data"].(map[string]interface{})
		fmt.Println("你回复的消息是: " + dataMap["message"].(string))
		if strings.Contains(dataMap["message"].(string), "[CQ:image") {
			sentGroupMessage(groupID, "福瑞添加成功: "+dataMap["message"].(string))
			msgDB = append(msgDB, dataMap["message"].(string))
		}
	}
	if strings.Contains(msgMap["message"].(string), "来只福瑞") {
		rand.Seed(time.Now().UnixNano())
		sentGroupMessage(groupID, msgDB[rand.Intn(len(msgDB))])
	}
}

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

	var groupMessage messages.GroupMessage
	json.Unmarshal(body, &groupMessage)
	log.Println("group message : ", groupMessage)

	log.Printf("收到一条消息：group_id = %d \n消息内容：%s",
		groupMessage.GroupId, groupMessage.Message)

	//if msgMap["message"] != nil && msgMap["group_id"] == 8.72367993e+08 {
	//
	//	groupManage(groupID, msgMap)
	//}
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
