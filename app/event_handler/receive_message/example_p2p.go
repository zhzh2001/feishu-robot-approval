package receiveMessage

import (
	"encoding/json"
	"xlab-feishu-robot/config"
	"xlab-feishu-robot/pkg/global"

	_ "github.com/sirupsen/logrus"
)

func init() {
	p2pMessageRegister(p2pHelpMenu, "help")
	// p2pMessageRegister(p2pTest, "test")
}

func p2pHelpMenu(messageevent *MessageEvent) {
	global.Cli.MessageSend("open_id", messageevent.Sender.Sender_id.Open_id, "text", "this is a P2P test string")
}

func p2pTest(messageevent *MessageEvent) {
	query := make(map[string]string)
	query["approval_code"] = config.C.Token.ApprovalCode
	query["start_time"] = "0"
	query["end_time"] = "1710204691000"
	resp := global.Cli.Request("get", "open-apis/approval/v4/instances", query, nil, nil)
	list := resp["instance_code_list"].([]interface{})
	for _, v := range list {
		resp := global.Cli.Request("get", "open-apis/approval/v4/instances/"+v.(string), nil, nil, nil)
		if resp["status"].(string) == "APPROVED" {
			var form []map[string]interface{}
			json.Unmarshal([]byte(resp["form"].(string)), &form)
			for _, v := range form {
				if v["name"].(string) == "采购事由" {
					// global.Cli.Send("open_id", messageevent.Sender.Sender_id.Open_id, "text", v["value"].(string))
				}
			}
		}
	}
}
