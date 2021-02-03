package watch

import "file-watch/message"

type Files struct {
    Files []string
}

func (f *Files) Change(msg chan message.Message) {
    if len(f.Files) == 0 {
        msg <- message.NewMessage("监测文件为空, 跳过", message.Debug)
    }
    
    
}
