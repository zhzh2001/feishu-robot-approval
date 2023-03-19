package eventHandler

import (
	"xlab-feishu-robot/app/event_handler/approval"
	receiveMessage "xlab-feishu-robot/app/event_handler/receive_message"
	"xlab-feishu-robot/pkg/dispatcher"
)

func Init() {
	// register your handlers here
	// example
	dispatcher.RegisterListener(receiveMessage.Receive, "im.message.receive_v1")
	dispatcher.RegisterListener(approval.Receive, "event_callback")
}
