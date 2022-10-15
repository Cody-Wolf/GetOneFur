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
	msgDB []string
}

func (rp *RandomPicture) Response(message messages.Messager) {
	groupMsg, ok := message.(*messages.GroupMessage)
	if ok == false {
		return
	}

	if strings.Contains(groupMsg.Message, "福瑞浓度") {
		sender.SendGroupMessage(strconv.FormatInt(groupMsg.GroupId, 10), "当前图库一共有 "+strconv.Itoa(len(rp.msgDB))+" 只福瑞")
	}

	if strings.Contains(groupMsg.Message, "[CQ:reply") &&
		strings.Contains(groupMsg.Message, "添加福瑞") {
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
			sender.SendGroupMessage(strconv.FormatInt(groupMsg.GroupId, 10), "福瑞添加成功: "+dataMap["message"].(string))
			rp.msgDB = append(rp.msgDB, dataMap["message"].(string))
		}
	}

	if strings.Contains(groupMsg.Message, "来只福瑞") {
		sender.SendGroupMessage(strconv.FormatInt(groupMsg.GroupId, 10), rp.msgDB[rand.Intn(len(rp.msgDB))])
	}
}

func init() {
	randSeed := time.Now().UnixNano()
	rand.Seed(randSeed)
	log.Println("rand seed = ", randSeed)
}
