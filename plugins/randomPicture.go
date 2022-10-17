package plugins

import (
	"GetOneFur/messages"
	"GetOneFur/sender"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type RandomPicture struct {
	msgDB         []string
	topic2Indexes map[string]map[int]bool
}

func (r *RandomPicture) GetPluginName() string {
	return "随机圣经"
}

func (r *RandomPicture) HelpInfo() string {
	return "/添加福瑞 (主题)：回复想要储存的消息，并只回复 “/添加福瑞”，即可保存，主题可选。\n" +
		"/来只福瑞 (主题)：从储存的消息中随机挑选一个发出。指定主题后只会有主题相关的图片。\n" +
		"/删除福瑞：回复想要删除的消息，并只回复 “/删除福瑞”，必须和已经存在的消息一模一样。\n" +
		"/福瑞浓度 (主题)：告诉你一共存了多少黑历史。"
}

func (r *RandomPicture) Response(message messages.Messager) {
	if r.topic2Indexes == nil {
		r.topic2Indexes = make(map[string]map[int]bool)
	}
	if r.topic2Indexes[""] == nil {
		r.topic2Indexes[""] = make(map[int]bool)
	}

	groupMsg, ok := message.(*messages.GroupMessage)
	if ok == false {
		return
	}

	words := strings.Split(groupMsg.Message, " ")
	topic := words[len(words)-1]
	if strings.HasPrefix(topic, "/") {
		topic = ""
	}
	log.Println("topic = ", topic)
	if r.topic2Indexes[topic] == nil {
		r.topic2Indexes[topic] = make(map[int]bool)
	}

	if strings.Count(groupMsg.Message, "/添加福瑞") > 1 {
		sender.SendGroupMessage(strconv.FormatInt(groupMsg.GroupId, 10), "您搁着套娃呢")
		return
	}

	if strings.HasPrefix(groupMsg.Message, "/福瑞浓度") {
		sender.SendGroupMessage(strconv.FormatInt(groupMsg.GroupId, 10), "当前图库一共有 "+strconv.Itoa(len(r.topic2Indexes[topic]))+" 只福瑞")
	}

	if strings.Contains(groupMsg.Message, "[CQ:reply") &&
		strings.Contains(groupMsg.Message, "/添加福瑞") {
		//sender.sendGroupMessage(groupID, "检测到回复")
		msgID := strings.TrimPrefix(groupMsg.Message, "[CQ:reply,id=")
		msgID, _, _ = strings.Cut(msgID, "]")
		fmt.Println("回复消息 ID : ", msgID)
		resBody := sender.SendMessage("get_msg", url.Values{"message_id": {msgID}})
		resMap := make(map[string]interface{})
		if err := json.Unmarshal(resBody, &resMap); err != nil {
			fmt.Printf("反序列化错误 err=%v\n", err)
			return
		}
		dataMap := resMap["data"].(map[string]interface{})
		fmt.Println("你回复的消息是: " + dataMap["message"].(string))
		if strings.Contains(dataMap["message"].(string), "[CQ:image") {
			sender.SendGroupMessage(strconv.FormatInt(groupMsg.GroupId, 10), topic+"福瑞添加成功: "+dataMap["message"].(string))
			r.msgDB = append(r.msgDB, dataMap["message"].(string))
			r.topic2Indexes[topic][len(r.msgDB)-1] = true
			r.topic2Indexes[""][len(r.msgDB)-1] = true
		}
	}

	if strings.Contains(groupMsg.Message, "[CQ:reply") &&
		strings.Contains(groupMsg.Message, "/删除福瑞") {
		//sender.sendGroupMessage(groupID, "检测到回复")
		msgID := strings.TrimPrefix(groupMsg.Message, "[CQ:reply,id=")
		msgID, _, _ = strings.Cut(msgID, "]")
		fmt.Println("回复消息 ID : ", msgID)
		resBody := sender.SendMessage("get_msg", url.Values{"message_id": {msgID}})
		resMap := make(map[string]interface{})
		if err := json.Unmarshal(resBody, &resMap); err != nil {
			fmt.Printf("反序列化错误 err=%v\n", err)
			return
		}
		dataMap := resMap["data"].(map[string]interface{})
		fmt.Println("你回复的消息是: " + dataMap["message"].(string))
		if strings.Contains(dataMap["message"].(string), "[CQ:image") {
			deleteFlag := false
			for index, msg := range r.msgDB {
				if strings.Contains(msg, dataMap["message"].(string)) {
					for name, indexes := range r.topic2Indexes {
						delete(indexes, index)
						log.Println("删除福瑞，topic = ", name)
					}
					deleteFlag = true
					break
				}
			}

			if deleteFlag {
				sender.SendGroupMessage(strconv.FormatInt(groupMsg.GroupId, 10), "福瑞删除成功: "+dataMap["message"].(string))
			} else {
				sender.SendGroupMessage(strconv.FormatInt(groupMsg.GroupId, 10), "不存在这个图片，福瑞删除失败: "+dataMap["message"].(string))
			}
		}
	}

	if strings.HasPrefix(groupMsg.Message, "/来只福瑞") {
		indexes := r.topic2Indexes[topic]
		randIndex := rand.Intn(len(indexes))
		count := 0
		for index := range indexes {
			if count == randIndex {
				sender.SendGroupMessage(strconv.FormatInt(groupMsg.GroupId, 10), r.msgDB[index])
				break
			}
			count++
		}
	}
}

func init() {
	randSeed := time.Now().UnixNano()
	rand.Seed(randSeed)
	log.Println("rand seed = ", randSeed)
}
