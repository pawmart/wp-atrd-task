package app

import (
	"github.com/gin-gonic/gin"
	"github.com/pawmart/wp-atrd-task/internal/storage"
)

type App struct {
	storage.Storage
}

//NewApp returns new app instance
func NewApp(s storage.Storage) *App {
	return &App{s}
}

//NewRoutes returns registered routes.
func (a *App) NewRoutes() *gin.Engine {
	r := gin.New()
	v1Router := r.Group("/v1")

	v1Router.POST("/secret", a.CreateSecretHandler())
	v1Router.GET("/secret/:id", a.GetSecretHandler())
	return r
}
