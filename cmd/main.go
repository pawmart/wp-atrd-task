package main

import (
	"github.com/alkmc/wp-atrd-task/internal/controller"
	"github.com/alkmc/wp-atrd-task/internal/repository"
	"github.com/alkmc/wp-atrd-task/internal/router"
	"github.com/alkmc/wp-atrd-task/internal/service"
	"github.com/alkmc/wp-atrd-task/internal/validator"
)

var (
	secretRepository = repository.NewRedis()
	secretService    = service.NewService(secretRepository)
	secretValidator  = validator.NewValidator()
	secretController = controller.NewController(secretService, secretValidator)
	secretRouter     = router.NewChiRouter()
)

func main() {
	mapUrls()
	defer secretRepository.CloseDB()
	secretRouter.SERVE()
}

func mapUrls() {
	secretRouter.GET("/v1/secret/{hash}", secretController.GetSecretByHash)
	secretRouter.POST("/v1/secret", secretController.AddSecret)
}
