package gateway

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler[T, M any] interface {
	Handle(*gin.Context, *T) (*M, error)
}

func handle[T, M any](handler Handler[T, M]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		resp, err := handler.Handle(c, &req)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
