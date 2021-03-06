package watch

import (
    "bytes"
    "file-watch/message"
    "fmt"
    "github.com/pkg/errors"
    "net/http"
    "os/exec"
    "time"
)

type Ngrok struct {
    Name     string
    FilePath string
}

func (n *Ngrok) Change(c chan message.Message) {
    if n.Name == "" {
        n.Name = "ngrok"
    }
    
    if n.FilePath == "" {
        c <- message.NewMessage(n.Name + " 日志文件未找到, 跳过", message.Debug)
        return
    }
    
    for {
        if n.ping() {
            if s, _ := n.status(); s != "0\n" {
                continue
            }
            // 通知
            c <- message.NewMessage("准备启动" + n.Name, message.Debug | message.Info)
            if err := exec.Command("bash", "-c", "systemctl start ngrok").Run(); err != nil {
                c <- message.NewMessage(fmt.Sprintf("启动 " + n.Name + " 失败, 错误信息: %v, 正在尝试重启...", err), message.Debug | message.Error)
                
                if err := exec.Command("bash", "-c", "systemctl restart ngrok").Run(); err != nil {
                    c <- message.NewMessage(fmt.Sprintf("重启 " + n.Name + " 失败, 错误信息: %v, 正在尝试重启...", err), message.Debug | message.Error)
                    close(c)
                }
            }

            if newChannel := n.newAddress(); newChannel != "" {
                c <- message.NewMessage(newChannel, message.Debug | message.Info)
            }
            
        } else {
            c <- message.NewMessage("网络连接不可达, 稍后重试", message.Debug)
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

func (n *Ngrok) status() (string, error) {
    var out bytes.Buffer
    cmd := exec.Command("/bin/bash", "-c", "systemctl status ngrok | grep running | wc -l")
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        return "", errors.Errorf("获取 ngrok 运行状态失败, 错误信息: %v", err)
    }
    
    return out.String(), nil
}

func (n *Ngrok) newAddress() string {
    var out bytes.Buffer
    cmd := exec.Command("/bin/bash", "-c", "grep 'started tunnel' " + n.FilePath + " | tail -1")
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        return fmt.Sprintf("获取 ngrok 最新地址失败, 错误信息: %v", err)
    }
    
    return out.String()
}
