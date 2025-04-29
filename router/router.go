package router

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/server"
)

//go:embed view
var FS embed.FS

// SetupRouter initializes and returns the Gin router.
func SetupRouter() *gin.Engine {
	r := gin.Default()

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
	frontendFS, err := fs.Sub(FS, "view")
	if err != nil {
		log.Fatalf("failed to initialize static files: %v", err)
	}
	r.StaticFS("/static", http.FS(frontendFS))
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

	r.LoadHTMLFiles("router/view/index.html", "router/view/channel_info.html")
	// if err := LoadHTMLFromEmbedFS(r, FS, "handler/view/index.html"); err != nil {
	// 	log.Fatalf("failed to load HTML templates: %v", err)
	// }
}
