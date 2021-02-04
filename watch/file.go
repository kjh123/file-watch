package watch

import (
    "file-watch/message"
    "github.com/fsnotify/fsnotify"
)

type Files struct {
    Files []string
}

func (f *Files) Change(msg chan message.Message) {
    if len(f.Files) == 0 {
        msg <- message.NewMessage("监测文件为空, 跳过", message.Debug)
        return
    }
    
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        msg <- message.NewMessage("文件监控初始化错误: " + err.Error(), message.Debug | message.Error)
        return
    }
    
    defer watcher.Close()
    
    for _, f := range f.Files {
        err = watcher.Add(f)
        if err != nil {
            msg <- message.NewMessage("文件监控错误: " + err.Error() + ", 文件名: " + f, message.Debug | message.Error)
            continue
        }
        msg <- message.NewMessage("添加监控文件: " + f, message.Debug)
    }
    
    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }
            
            msg <- message.NewMessage("文件发生变动: " + event.String(), message.Debug | message.Info)
        
        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            msg <- message.NewMessage("文件监控错误: " + err.Error(), message.Debug | message.Error)
        }
    }
}
