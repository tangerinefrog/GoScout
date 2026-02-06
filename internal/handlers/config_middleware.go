package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) ConfigMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config, err := h.configRepository.Get(c.Request.Context())
		if err != nil {
			msg := fmt.Sprintf("Could not load config: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		h.config = config

		c.Next()
	}
}
