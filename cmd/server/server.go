package main

import (
	"github.com/systemz/wp-atrd-task/internal/model"
	"github.com/systemz/wp-atrd-task/internal/rest"
)

func main() {
	// start connection to DB
	model.RedisInit()

	// start REST API
	rest.StartWebServer()
}
