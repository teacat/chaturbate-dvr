package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/adapter/manager"
	"github.com/teacat/chaturbate-dvr/gateway"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "chaturbate-dvr",
		Usage: "",
		Action: func(*cli.Context) error {
			r := gin.Default()
			manager, err := manager.New()
			if err != nil {
				return fmt.Errorf("new manager: %w", err)
			}
			g := gateway.New(manager)
			for route, handler := range g.Routes() {
				r.POST(route, handler)
			}
			return r.Run()
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
