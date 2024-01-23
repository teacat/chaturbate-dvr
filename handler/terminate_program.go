package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

//=======================================================
// Request & Response
//=======================================================

type TerminateProgramRequest struct {
}

type TerminateProgramResponse struct {
}

//=======================================================
// Factory
//=======================================================

type TerminateProgramHandler struct {
	cli *cli.Context
}

func NewTerminateProgramHandler(cli *cli.Context) *TerminateProgramHandler {
	return &TerminateProgramHandler{cli}
}

//=======================================================
// Handle
//=======================================================

func (h *TerminateProgramHandler) Handle(c *gin.Context) {
	var req *TerminateProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	log.Println("program terminated by user request, see ya 👋")
	os.Exit(0)
}
