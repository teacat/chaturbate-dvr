package handler

import "github.com/gin-gonic/gin"

//=======================================================
// Request & Response
//=======================================================

//=======================================================
// Factory
//=======================================================

type ViewIndexHandler struct {
}

func NewViewIndexHandler() *ViewIndexHandler {
	return &ViewIndexHandler{}
}

//=======================================================
// Handle
//=======================================================

func (h *ViewIndexHandler) Handle(ctx *gin.Context) {
	ctx.HTML(200, "index.tmpl", gin.H{})
}
