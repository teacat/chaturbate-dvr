package view

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

//go:embed templates
var FS embed.FS

// InfoTpl is a template for rendering channel information.
var InfoTpl *template.Template

func init() {
	var err error

	InfoTpl, err = template.New("update").ParseFS(FS, "templates/channel_info.html")
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}
}

// StaticFS initializes the static file system for serving frontend files.
func StaticFS() (http.FileSystem, error) {
	frontendFS, err := fs.Sub(FS, "templates")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize static files: %w", err)
	}
	return http.FS(frontendFS), nil
}
