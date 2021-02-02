package notify

import (
    "file-watch/message"
    "log"
)

type DingDing struct {
    Url   string
    Level message.MsgLevel
}

func (d *DingDing) NotifyLevel() message.MsgLevel {
    return d.Level
}

func (d *DingDing) Notify(msg message.Message) {
    log.Println("钉钉消息:", msg.String())
}
