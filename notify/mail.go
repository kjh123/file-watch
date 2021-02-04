package notify

import "file-watch/message"

type Email struct {
    Account  string
    Password string
    Level    message.MsgLevel
}

func (m *Email) Notify(msg message.Message) {

}

func (m *Email) NotifyLevel() message.MsgLevel {
    return m.Level
}

func (m *Email) HasError() bool {
    return false
}
