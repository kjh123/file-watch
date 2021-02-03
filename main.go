package main

import (
    "file-watch/message"
    "file-watch/notify"
    "file-watch/watch"
    "fmt"
    "github.com/urfave/cli/v2"
    "os"
)

const (
    version = "v1.1.0"
)

var (
    LogPath        string
    NgrokLogPath   string
    NotifyDingPath string
)

func main() {
    app := &cli.App{
        Name:                   "file-watch",
        Usage:                  "文件变更监控",
        UseShortOptionHandling: true,
        HideHelpCommand:        true,
        Version:                version,
        Flags:                  getArgs(),
        Action:                 watching(),
    }
    
    if err := app.Run(os.Args); err != nil {
        fmt.Println(err)
    }
}

func watching() cli.ActionFunc {
    return func(c *cli.Context) error {
        
        w := watch.FileWatch{
            Channels: []notify.Channel{
                &notify.DingDing{Url: NotifyDingPath, Level: message.Debug | message.Info},
                &notify.Log{LogPath: LogPath, Level: message.Debug | message.Info | message.Error},
                // TODO Email
                // &notify.Email{Account: "", Password: ""},
            },
            Files: []watch.File{
                &watch.Files{Files: c.Args().Slice()},
                &watch.Ngrok{FilePath: NgrokLogPath},
            },
        }
        w.Run()
        return nil
    }
}

// getArgs run command with arguments
func getArgs() []cli.Flag {
    return []cli.Flag{
        &cli.StringFlag{
            Name:        "ding_url",
            Destination: &NotifyDingPath,
            EnvVars:     []string{"DING_URL"},
            Usage:       "钉钉消息通知",
        },
        &cli.StringFlag{
            Name:        "log",
            Destination: &LogPath,
            EnvVars:     []string{"LOG"},
            Value:       "/var/log/file_watch.log",
            Usage:       "是否记录文件变动到日志",
        },
        &cli.StringFlag{
            Name:        "ngrok_log_path",
            EnvVars:     []string{"NGROK_LOG_PATH"},
            Destination: &NgrokLogPath,
            Usage:       "Ngrok 日志文件路径",
        },
    }
}
