package watch

import (
    "net/http"
    "time"
)

type Ngrok struct {
    Name     string
    FilePath string
}

func (n *Ngrok) Change(c chan string) {
    if n.Name == "" {
        n.Name = "ngrok"
    }
    
    if n.FilePath == "" {
        c <- n.Name + " 文件未找到"
        return
    }
    
    for {
        if n.Ping() {
            // TODO
        } else {
            time.Sleep(time.Minute)
        }
    }
}

func (n *Ngrok) Ping() bool {
    res, _ := http.Head("https://www.baidu.com")
    if res != nil && res.StatusCode == 200 {
        defer res.Body.Close()
        return true
    }
    
    return false
}
