package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/teacat/chaturbate-dvr/chaturbate"
	"github.com/urfave/cli/v2"
)

type UpdateLogLevelHandler struct {
	cli *cli.Context
}

// Custom validator for LogType
func LogTypeValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	switch value {
	case string(chaturbate.LogTypeDebug), string(chaturbate.LogTypeInfo), string(chaturbate.LogTypeWarning), string(chaturbate.LogTypeError):
		return true
	}
	return false
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("logtype", LogTypeValidator)
	}
}

func NewUpdateLogLevelHandler(cli *cli.Context) *UpdateLogLevelHandler {
	return &UpdateLogLevelHandler{cli}
}

func (h *UpdateLogLevelHandler) Handle(c *gin.Context) {
	var req chaturbate.LogLevelRequest

	// Bind and validate the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format. Expected {\"log_level\": \"INFO\"}",
		})
		return
	}

	// Use the correct log type for setting the global log level
	chaturbate.SetGlobalLogLevel(req.LogLevel)

	log.Printf("Global log level updated to: %s", req.LogLevel)

	// Send success response
	c.JSON(http.StatusOK, gin.H{
		"message":   "Log level updated",
		"log_level": req.LogLevel,
	})
}

// func (h *UpdateLogLevelHandler) Handle(c *gin.Context) {
// 	// Read the raw request body for debugging
// 	bodyBytes, err := c.GetRawData()
// 	if err != nil {
// 		log.Printf("Error reading request body: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
// 		return
// 	}

// 	// Log the raw request body
// 	log.Printf("Received raw request body: %s", string(bodyBytes))

// 	// Reset the request body so it can be re-read by ShouldBindJSON
// 	c.Request.Body = ioutil.NopCloser(strings.NewReader(string(bodyBytes)))

// 	// Attempt to bind the JSON to the struct
// 	var req LogLevelRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		log.Printf("Error binding JSON: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid request format. Expected {\"log_level\": \"INFO\"}",
// 		})
// 		return
// 	}

// 	// Log the updated log level
// 	log.Printf("Log level updated to: %s", req.LogLevel)

// 	// Store the log level in the CLI context if needed
// 	h.cli.Set("log_level", string(req.LogLevel))

// 	// Send success response
// 	c.JSON(http.StatusOK, gin.H{
// 		"message":   "Log level updated",
// 		"log_level": req.LogLevel,
// 	})
// }

// NewUpdateLogLevelHandler creates a handler for updating log level.
// func NewUpdateLogLevelHandler(c *cli.Context) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		var req LogLevelRequest

// 		// Bind and validate request body
// 		if err := ctx.ShouldBindJSON(&req); err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
// 			return
// 		}

// 		if !allowedLogLevels[req.LogLevel] {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log level"})
// 			return
// 		}

// 		ctx.JSON(http.StatusOK, gin.H{
// 			"message":   "Log level updated",
// 			"log_level": req.LogLevel,
// 		})
// 	}
// }
