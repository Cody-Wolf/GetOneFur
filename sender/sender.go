package sender

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func SendMessage(opt string, msgMap url.Values) []byte {
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

func SendGroupMessage(groupID, message string) []byte {
	return SendMessage("send_group_msg", url.Values{"group_id": {groupID}, "message": {message}})
}
