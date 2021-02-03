package notify

import (
    "file-watch/message"
    "log"
    "os"
)

type Log struct {
    LogPath string
    Level   message.MsgLevel
}

func (l *Log) NotifyLevel() message.MsgLevel {
    return l.Level
}

func (l *Log) Notify(msg message.Message) {
    if len(l.LogPath) > 0 {
        f, err := os.OpenFile(l.LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
        if err != nil {
            l.err(err.Error())
            return
        }
        defer f.Close()
    
        l := log.New(f, "", log.LstdFlags)
        l.Println(msg.String())
    }
    
}

func (l *Log) err (str string) {
    log.Println("log error: " + str)
}
