package watch

import (
    "file-watch/notify"
    "log"
    "os"
)

type File interface {
    Change(chan string)
}

type FileWatch struct {
    Channels   []notify.Channel
    Files      []File
    LogEnabled bool
    LogPath    string
}

func (w *FileWatch) HasChanged() bool {
    return false
}

func (w *FileWatch) SetLogEnabled(enabled bool) {
    w.LogEnabled = enabled
}

func (w *FileWatch) AppendChannel(c notify.Channel) {
    w.Channels = append(w.Channels, c)
}

func (w *FileWatch) Run() {
    if len(w.Files) == 0 {
        return
    }
    
    c := make(chan string)
    for _, f := range w.Files {
        go f.Change(c)
    }
    defer close(c)
    
    for {
        select {
        case str := <-c:
            w.Logging(str)
        }
    }
}

func (w *FileWatch) Logging(str string) {
    if len(w.Channels) > 0 {
        for _, c := range w.Channels {
            c.Notify()
        }
    }
    
    if w.LogEnabled {
        f, err := os.OpenFile(w.LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
        if err != nil {
            log.Println(err)
        }
        defer f.Close()
        
        f.WriteString(str)
    }
}
