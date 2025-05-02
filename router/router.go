package router

import (
	"embed"
	"html/template"
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/router/view"
	"github.com/teacat/chaturbate-dvr/server"
)

// SetupRouter initializes and returns the Gin router.
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	if err := LoadHTMLFromEmbedFS(r, view.FS, "templates/index.html", "templates/channel_info.html"); err != nil {
		log.Fatalf("failed to load HTML templates: %v", err)
	}

	// Apply authentication if configured
	SetupAuth(r)
	// Serve static frontend files
	SetupStatic(r)
	// Register views
	SetupViews(r)

	return r
}

// SetupAuth applies basic authentication if credentials are provided.
func SetupAuth(r *gin.Engine) {
	if server.Config.AdminUsername != "" && server.Config.AdminPassword != "" {
		auth := gin.BasicAuth(gin.Accounts{
			server.Config.AdminUsername: server.Config.AdminPassword,
		})
		r.Use(auth)
	}
}

// SetupStatic serves static frontend files.
func SetupStatic(r *gin.Engine) {
	fs, err := view.StaticFS()
	if err != nil {
		log.Fatalf("failed to initialize static files: %v", err)
	}
	r.StaticFS("/static", fs)
}

// setupViews registers HTML templates and view handlers.
func SetupViews(r *gin.Engine) {
	r.GET("/", Index)
	r.GET("/updates", Updates)
	r.POST("/update_config", UpdateConfig)
	r.POST("/create_channel", CreateChannel)
	r.POST("/stop_channel/:username", StopChannel)
	r.POST("/pause_channel/:username", PauseChannel)
	r.POST("/resume_channel/:username", ResumeChannel)

}

// LoadHTMLFromEmbedFS loads specific HTML templates from an embedded filesystem and registers them with Gin.
func LoadHTMLFromEmbedFS(r *gin.Engine, embeddedFS embed.FS, files ...string) error {
	templ := template.New("")
	for _, file := range files {
		content, err := embeddedFS.ReadFile(file)
		if err != nil {
			return err
		}
		_, err = templ.New(filepath.Base(file)).Parse(string(content))
		if err != nil {
			return err
		}
	}

	// Set the parsed templates as the HTML renderer for Gin
	r.SetHTMLTemplate(templ)
	return nil
}
