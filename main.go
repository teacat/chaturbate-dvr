package main

import (
	"fmt"
	"log"
	"os"

	"github.com/teacat/chaturbate-dvr/config"
	"github.com/teacat/chaturbate-dvr/entity"
	"github.com/teacat/chaturbate-dvr/manager"
	"github.com/teacat/chaturbate-dvr/router"
	"github.com/teacat/chaturbate-dvr/server"
	"github.com/urfave/cli/v2"
)

const logo = `
 ██████╗██╗  ██╗ █████╗ ████████╗██╗   ██╗██████╗ ██████╗  █████╗ ████████╗███████╗
██╔════╝██║  ██║██╔══██╗╚══██╔══╝██║   ██║██╔══██╗██╔══██╗██╔══██╗╚══██╔══╝██╔════╝
██║     ███████║███████║   ██║   ██║   ██║██████╔╝██████╔╝███████║   ██║   █████╗
██║     ██╔══██║██╔══██║   ██║   ██║   ██║██╔══██╗██╔══██╗██╔══██║   ██║   ██╔══╝
╚██████╗██║  ██║██║  ██║   ██║   ╚██████╔╝██║  ██║██████╔╝██║  ██║   ██║   ███████╗
 ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚══════╝
██████╗ ██╗   ██╗██████╗
██╔══██╗██║   ██║██╔══██╗
██║  ██║██║   ██║██████╔╝
██║  ██║╚██╗ ██╔╝██╔══██╗
██████╔╝ ╚████╔╝ ██║  ██║
╚═════╝   ╚═══╝  ╚═╝  ╚═╝`

func main() {
	app := &cli.App{
		Name:    "chaturbate-dvr",
		Version: "2.0.3",
		Usage:   "Record your favorite Chaturbate streams automatically. 😎🫵",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Usage:   "The username of the channel to record",
				Value:   "",
			},
			&cli.StringFlag{
				Name:  "admin-username",
				Usage: "Username for web authentication (optional)",
				Value: "",
			},
			&cli.StringFlag{
				Name:  "admin-password",
				Usage: "Password for web authentication (optional)",
				Value: "",
			},
			&cli.IntFlag{
				Name:  "framerate",
				Usage: "Desired framerate (FPS)",
				Value: 30,
			},
			&cli.IntFlag{
				Name:  "resolution",
				Usage: "Desired resolution (e.g., 1080 for 1080p)",
				Value: 1080,
			},
			&cli.StringFlag{
				Name:  "pattern",
				Usage: "Template for naming recorded videos",
				Value: "videos/{{.Username}}_{{.Year}}-{{.Month}}-{{.Day}}_{{.Hour}}-{{.Minute}}-{{.Second}}{{if .Sequence}}_{{.Sequence}}{{end}}",
			},
			&cli.IntFlag{
				Name:  "max-duration",
				Usage: "Split video into segments every N minutes ('0' to disable)",
				Value: 0,
			},
			&cli.IntFlag{
				Name:  "max-filesize",
				Usage: "Split video into segments every N MB ('0' to disable)",
				Value: 0,
			},
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "Port for the web interface and API",
				Value:   "8080",
			},
			&cli.IntFlag{
				Name:  "interval",
				Usage: "Check if the channel is online every N minutes",
				Value: 1,
			},
			&cli.StringFlag{
				Name:  "cookies",
				Usage: "Cookies to use in the request (format: key=value; key2=value2)",
				Value: "",
			},
			&cli.StringFlag{
				Name:  "user-agent",
				Usage: "Custom User-Agent for the request",
				Value: "",
			},
			&cli.StringFlag{
				Name:  "domain",
				Usage: "Chaturbate domain to use",
				Value: "https://chaturbate.com/",
			},
		},
		Action: start,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func start(c *cli.Context) error {
	fmt.Println(logo)

	var err error
	server.Config, err = config.New(c)
	if err != nil {
		return fmt.Errorf("new config: %w", err)
	}
	server.Manager, err = manager.New()
	if err != nil {
		return fmt.Errorf("new manager: %w", err)
	}

	// init web interface if username is not provided
	if server.Config.Username == "" {
		fmt.Printf("👋 Visit http://localhost:%s to use the Web UI\n\n\n", c.String("port"))

		if err := server.Manager.LoadConfig(); err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		return router.SetupRouter().Run(":" + c.String("port"))
	}

	// else create a channel with the provided username
	if err := server.Manager.CreateChannel(&entity.ChannelConfig{
		IsPaused:    false,
		Username:    c.String("username"),
		Framerate:   c.Int("framerate"),
		Resolution:  c.Int("resolution"),
		Pattern:     c.String("pattern"),
		MaxDuration: c.Int("max-duration"),
		MaxFilesize: c.Int("max-filesize"),
	}, false); err != nil {
		return fmt.Errorf("create channel: %w", err)
	}

	// block forever
	select {}
}
