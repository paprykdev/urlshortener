package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/paprykdev/urlshortener/internal/handlers"
)

func LinksRoutes (r *gin.Engine, handler *handlers.LinksHandler) {
	links := r.Group("/links")

	links.GET("", handler.GetAll)
	links.POST("", handler.CreateLink)
}
