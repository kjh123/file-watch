package watch

import (
    "file-watch/message"
    "file-watch/notify"
    "log"
    "os"
    "os/signal"
)

type File interface {
    Change(chan message.Message)
}

type FileWatch struct {
    Channels   []notify.Channel
    Files      []File
}

func NewFileWatch() *FileWatch {
    return &FileWatch{}
}

func (w *FileWatch) AppendFiles(f ...File) {
    w.Files = append(w.Files, f...)
}

func (w *FileWatch) AppendChannel(c ...notify.Channel) {
    w.Channels = append(w.Channels, c...)
}

func (w *FileWatch) Run() {
    if len(w.Files) == 0 {
        log.Println("watch files is empty")
        return
    }
    
    c := make(chan message.Message)
    go w.logging(c)
    for _, f := range w.Files {
        go f.Change(c)
    }
    defer close(c)
    
    quit := make(chan os.Signal)
    signal.Notify(quit, os.Interrupt)
    <-quit
}

// logging WriteString message.
func (w *FileWatch) logging(c chan message.Message) {
    for {
        select {
        case msg, ok := <-c:
            if ok {
                log.Println(msg.String())
        
                for _, c := range w.Channels {
                    if !c.HasError() && c.NotifyLevel()&msg.Level != 0 {
                        go c.Notify(msg)
                    }
                }
            }
            
        }
    }
}
