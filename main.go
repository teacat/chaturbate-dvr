package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/chaturbate"
	"github.com/teacat/chaturbate-dvr/handler"
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
		Version: "1.0.7",
		Usage:   "Records your favorite Chaturbate stream 😎🫵",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Usage:   "channel username to record",
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "gui-username",
				Aliases: []string{"gui-u"},
				Usage:   "username for auth web (optional)",
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "gui-password",
				Aliases: []string{"gui-p"},
				Usage:   "password for auth web (optional)",
				Value:   "",
			},
			&cli.IntFlag{
				Name:    "framerate",
				Aliases: []string{"f"},
				Usage:   "preferred framerate",
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
				Usage: "log level, available: 'DEBUG', 'INFO', 'WARN', 'ERROR'",
				Value: "INFO",
			},
			&cli.StringFlag{
				Name:  "port",
				Usage: "port to expose the web interface and API",
				Value: "8080",
			},
			&cli.IntFlag{
				Name:    "interval",
				Aliases: []string{"i"},
				Usage:   "minutes to check if the channel is online",
				Value:   1,
			},
			&cli.StringFlag{
				Name:	"cf-cookie",
				Usage:   "Cloudflare cookie to bypass anti-bot page",
				Value:   "",
			},
			&cli.StringFlag{
				Name:  "user-agent",
				Usage: "Custom user agent for when using cf-cookie",
				Value: "",
			},
			//&cli.StringFlag{
			//	Name:  "gui",
			//	Usage: "enabling GUI, availables: 'no', 'web'",
			//	Value: "web",
			//},
		},
		Action: start,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func start(c *cli.Context) error {
	fmt.Println(logo)
	if c.String("cf-cookie") != "" && c.String("user-agent") == ""{
		return fmt.Errorf("When using the cf-cookie option a user-agent MUST be supplied")
	}

	//if c.String("gui") == "web" {
	if c.String("username") == "" {
		return startWeb(c)
	}

	m := chaturbate.NewManager(c)
	if err := m.CreateChannel(&chaturbate.Config{
		Username:           c.String("username"),
		Framerate:          c.Int("framerate"),
		Resolution:         c.Int("resolution"),
		ResolutionFallback: c.String("resolution-fallback"),
		FilenamePattern:    c.String("filename-pattern"),
		SplitDuration:      c.Int("split-duration"),
		SplitFilesize:      c.Int("split-filesize"),
		Interval:           c.Int("interval"),
		CFCookie:			c.String("cf-cookie"),
		UserAgent:			c.String("user-agent"),
	}); err != nil {
		return err
	}

	select {} // block forever
}

//go:embed handler/view
var FS embed.FS

func startWeb(c *cli.Context) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	//r.Use(cors.Default())
	m := chaturbate.NewManager(c)
	if err := m.LoadChannels(); err != nil {
		return err
	}

	fe, err := fs.Sub(FS, "handler/view")
	if err != nil {
		log.Fatalln(err)
	}
	guiUsername := c.String("gui-username")
	guiPassword := c.String("gui-password")

	logLevel := c.String("log-level")

	chaturbate.InitGlobalLogLevel(chaturbate.LogType(logLevel))

	var authorized = r.Group("/")
	var authorizedApi = r.Group("/api")

	if guiUsername != "" && guiPassword != "" {
		ginBasicAuth := gin.BasicAuth(gin.Accounts{
			guiUsername: guiPassword,
		})
		authorized.Use(ginBasicAuth)
		authorizedApi.Use(ginBasicAuth)
	}

	authorized.StaticFS("/static", http.FS(fe))
	authorized.StaticFileFS("/", "/", http.FS(fe))

	authorizedApi.POST("/get_channel", handler.NewGetChannelHandler(m, c).Handle)
	authorizedApi.POST("/create_channel", handler.NewCreateChannelHandler(m, c).Handle)
	authorizedApi.POST("/list_channels", handler.NewListChannelsHandler(m, c).Handle)
	authorizedApi.POST("/delete_channel", handler.NewDeleteChannelHandler(m, c).Handle)
	authorizedApi.POST("/pause_channel", handler.NewPauseChannelHandler(m, c).Handle)
	authorizedApi.POST("/resume_channel", handler.NewResumeChannelHandler(m, c).Handle)
	authorizedApi.GET("/listen_update", handler.NewListenUpdateHandler(m, c).Handle)
	authorizedApi.POST("/get_settings", handler.NewGetSettingsHandler(c).Handle)
	authorizedApi.POST("/terminate_program", handler.NewTerminateProgramHandler(c).Handle)
	authorizedApi.POST("/update_log_level", handler.NewUpdateLogLevelHandler(c).Handle)

	fmt.Printf("👋 Visit http://localhost:%s to use the Web UI\n", c.String("port"))

	return r.Run(fmt.Sprintf(":%s", c.String("port")))
}
