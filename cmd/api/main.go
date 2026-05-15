package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paprykdev/urlshortener/internal/db"
	"github.com/paprykdev/urlshortener/internal/handlers"
	"github.com/paprykdev/urlshortener/internal/routes"
)

func main() {
	r := gin.Default()
	db := db.Init()

	defer db.Close()

	linksHandler := handlers.NewLinksHandler(db)

	routes.LinksRoutes(r, linksHandler)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET(":code", linksHandler.Redirect)

	r.Run()
}
