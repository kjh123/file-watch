package main

import (
    "file-watch/notify"
    "file-watch/watch"
    "fmt"
    "github.com/urfave/cli/v2"
    "os"
)

const (
    version = "v1.0.0"
    logPath = "/var/log/file-watch.log"
)

var (
    NgrokLogPath   string
    LogEnabled     bool
    NotifyDingPath string
)

func main() {
    app := &cli.App{
        Name:                   "file-watch",
        Usage:                  "文件变更监控",
        UseShortOptionHandling: true,
        HideHelpCommand:        true,
        Version:                version,
        Flags: []cli.Flag{
            &cli.BoolFlag{
                Name:        "log",
                Destination: &LogEnabled,
                Value:       true,
                Usage:       "是否记录文件变动到日志",
            },
            &cli.StringFlag{
                Name:        "ding_url",
                Destination: &NotifyDingPath,
                EnvVars:     []string{"DING_URL"},
                Usage:       "钉钉消息通知",
            },
            &cli.StringFlag{
                Name:        "ngrok_log_path",
                EnvVars:     []string{"NGROK_LOG_PATH"},
                Destination: &NgrokLogPath,
                Usage:       "Ngrok 日志文件路径",
            },
        },
        Action: func(c *cli.Context) error {
            w := watch.FileWatch{
                Channels: []notify.Channel{
                    &notify.DingDing{Url: NotifyDingPath},
                    // TODO Email
                    // &notify.Email{Account: "", Password: ""},
                },
                Files: []watch.File{
                    &watch.Ngrok{FilePath: NgrokLogPath},
                },
                LogEnabled: LogEnabled,
                LogPath:    logPath,
            }
            w.Run()
            return nil
        },
    }
    
    if err := app.Run(os.Args); err != nil {
        fmt.Println(err)
    }
}
