package notify

import (
    "bytes"
    "encoding/json"
    "file-watch/message"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "time"
)

var appName = "file-watch"

type DingDing struct {
    Url   string
    Level message.MsgLevel
}

func (d *DingDing) NotifyLevel() message.MsgLevel {
    return d.Level
}

func (d *DingDing) Notify(msg message.Message) {
    if u, err := url.Parse(d.Url); (u != nil && u.String() == "") || err != nil {
        d.err("please check ding notify url")
        return
    }
    
    markdown := ``
    markdown += fmt.Sprintf("## %s %s\n", appName, time.Now().Format("01-02 15:04:05"))
    markdown += fmt.Sprintf("### %s\n", msg.String())
    
    dm := dingMsg{Msgtype: "markdown"}
    dm.Markdown.Title = appName + " notify"
    dm.Markdown.Text = markdown
    
    bs, err := json.Marshal(dm)
    if err != nil {
        d.err(err.Error())
        return
    }
    
    res, e := http.Post(d.Url, "application/json", bytes.NewBuffer(bs))
    if e != nil {
        d.err(e.Error())
        return
    }
    defer res.Body.Close()
}

func (d *DingDing) err(str string) {
    log.Println("ding notify: " + str)
}

type dingMsg struct {
    Msgtype string `json:"msgtype"`
    Text    struct {
        Content string `json:"content"`
    } `json:"text"`
    Markdown struct {
        Title string `json:"title"`
        Text  string `json:"text"`
    } `json:"markdown"`
}
