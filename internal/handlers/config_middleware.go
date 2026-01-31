package handlers

import (
	"fmt"
	"net/http"

	"github.com/tangerinefrog/GoScout/internal/data/repositories"

	"github.com/gin-gonic/gin"
)

func (h *handler) ConfigMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		configRepo := repositories.NewConfigRepo(h.db)
		config, err := configRepo.Get(c.Request.Context())
		if err != nil {
			msg := fmt.Sprintf("Could not load config: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		h.config = config

		c.Next()
	}
}
