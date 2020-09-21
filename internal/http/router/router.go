package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pawmart/wp-atrd-task/internal/http/app"
	"net/http"
)

//NewRoutes returns registered routes.
func NewRoutes(a *app.App) http.Handler {
	r := gin.New()
	v1Router := r.Group("/v1")

	v1Router.POST("/secret", a.CreateSecretHandler())
	v1Router.GET("/secret/:id", a.GetSecretHandler())

	return r
}
