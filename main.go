package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/chaturbate"
	"github.com/teacat/chaturbate-dvr/handler"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "chaturbate-dvr",
		Usage: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Usage:   "channel username to record.",
				Value:   "",
			},
			&cli.IntFlag{
				Name:    "framerate",
				Aliases: []string{"f"},
				Usage:   "preferred framerate.",
				Value:   30,
			},
			&cli.IntFlag{
				Name:    "resolution",
				Aliases: []string{"r"},
				Usage:   "preferred resolution",
				Value:   1080,
			},
			&cli.StringFlag{
				Name:    "resolution-fallback",
				Aliases: []string{"rf"},
				Usage:   "fallback to 'up' (larger) or 'down' (smaller) resolution if preferred resolution is not available",
				Value:   "down",
			},
			&cli.StringFlag{
				Name:    "filename-pattern",
				Aliases: []string{"fp"},
				Usage:   "filename pattern for videos",
				Value:   "videos/{{.Username}}_{{.Year}}-{{.Month}}-{{.Day}}_{{.Hour}}-{{.Minute}}-{{.Second}}{{if .Sequence}}_{{.Sequence}}{{end}}",
			},
			&cli.IntFlag{
				Name:    "split-duration",
				Aliases: []string{"sd"},
				Usage:   "minutes to split each video into segments ('0' to disable)",
				Value:   0,
			},
			&cli.IntFlag{
				Name:    "split-filesize",
				Aliases: []string{"sf"},
				Usage:   "size in MB to split each video into segments ('0' to disable)",
				Value:   0,
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "log level, availables: 'debug', 'info', 'error'",
				Value: "info",
			},
			&cli.StringFlag{
				Name:  "port",
				Usage: "port to expose the web interface and API",
				Value: "8080",
			},
			&cli.StringFlag{
				Name:  "gui",
				Usage: "enabling GUI, availables: 'none', 'web'",
				Value: "none",
			},
		},
		Action: start,
		Commands: []*cli.Command{
			{
				Name:   "start",
				Usage:  "",
				Action: start,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func start(c *cli.Context) error {
	r := gin.Default()
	r.Use(cors.Default())
	m, err := chaturbate.NewManager()
	if err != nil {
		return err
	}
	r.POST("/api/get_channel", handler.NewGetChannelHandler(m, c).Handle)
	r.POST("/api/create_channel", handler.NewCreateChannelHandler(m, c).Handle)
	r.POST("/api/list_channels", handler.NewListChannelsHandler(m, c).Handle)
	r.POST("/api/delete_channel", handler.NewDeleteChannelHandler(m, c).Handle)
	r.POST("/api/pause_channel", handler.NewPauseChannelHandler(m, c).Handle)
	r.POST("/api/resume_channel", handler.NewResumeChannelHandler(m, c).Handle)
	r.GET("/api/listen_update", handler.NewListenUpdateHandler(m, c).Handle)
	r.POST("/api/get_settings", handler.NewGetSettingsHandler(c).Handle)

	return r.Run(fmt.Sprintf(":%s", c.String("port")))
}
