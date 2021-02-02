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
    LogPath    string
}

func (w *FileWatch) HasChanged() bool {
    return false
}

func (w *FileWatch) AppendChannel(c notify.Channel) {
    w.Channels = append(w.Channels, c)
}

func (w *FileWatch) Run() {
    if len(w.Files) == 0 {
        log.Println("watch files is empty")
        return
    }
    
    c := make(chan message.Message)
    go w.readChannel(c)
    for _, f := range w.Files {
        go f.Change(c)
    }
    defer close(c)
    
    quit := make(chan os.Signal)
    signal.Notify(quit, os.Interrupt)
    <-quit
}

func (w *FileWatch) readChannel(c chan message.Message) {
    for {
        select {
        case msg, ok := <-c:
            if ok {
                w.Logging(msg)
            }
            
        }
    }
}

func (w *FileWatch) Logging(msg message.Message) {
    // Stdout
    log.Println(msg.String())
    
    for _, c := range w.Channels {
        if c.NotifyLevel() & msg.Level != 0 {
            go c.Notify(msg)
        }
    }
    
    if len(w.LogPath) > 0 {
        f, err := os.OpenFile(w.LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
        if err != nil {
            log.Println(err)
        }
        defer f.Close()
        
        f.WriteString(msg.String())
    }
}
