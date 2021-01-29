package watch

import (
    "bytes"
    "log"
    "net/http"
    "os/exec"
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
        if n.ping() {
            if n.status() == "0" {
                if err := exec.Command("bash", "-c", "systemctl start ngrok").Run(); err != nil {
                    log.Printf("启动 ngrok 失败, 错误信息: %v, 正在尝试重启...", err)
                    if err := exec.Command("bash", "-c", "systemctl restart ngrok").Run(); err != nil {
                        log.Printf("重启 ngrok 失败, 错误信息: %v", err)
                        close(c)
                    }
                }
                if newChannel := n.newAddress(); newChannel != "" {
                    c <- n.newAddress()
                }
            }
        } else {
            time.Sleep(time.Minute)
        }
    }
}

func (n *Ngrok) ping() bool {
    res, _ := http.Head("https://www.baidu.com")
    if res != nil && res.StatusCode == 200 {
        defer res.Body.Close()
        return true
    }
    
    return false
}

func (n *Ngrok) status() string {
    var out bytes.Buffer
    cmd := exec.Command("/bin/bash", "-c", "systemctl status ngrok | grep running | wc -l")
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        log.Printf("获取 ngrok 运行状态失败, 错误信息: %v", err)
        return "0"
    }
    
    return out.String()
}

func (n *Ngrok) newAddress() string {
    var out bytes.Buffer
    cmd := exec.Command("/bin/bash", "-c", "grep 'started tunnel' /var/log/ngrok.log | tail -1")
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        log.Printf("获取 ngrok 最新地址失败, 错误信息: %v", err)
        return ""
    }
    
    return out.String()
}
