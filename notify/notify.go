package notify

import "file-watch/message"

type Channel interface {
    NotifyLevel() message.MsgLevel
    Notify(msg message.Message)
}
