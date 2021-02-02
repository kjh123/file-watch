package message

type MsgLevel uint8

const (
    Debug MsgLevel = 1 << iota
    Info
    Error
)

type Message struct {
    Level MsgLevel
    Msg   string
}

func NewMessage(msg string, level ...MsgLevel) Message {
    return Message{
        Level: level[0],
        Msg: msg,
    }
}

func (m *Message) String() string {
    return m.Msg
}

func (m *Message) SetLevel(ml MsgLevel) {
    m.Level = ml
}

func (m *Message) SetMsg(msg string) {
    m.Msg = msg
}
