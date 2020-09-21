package app

import "github.com/pawmart/wp-atrd-task/internal/storage"

type App struct {
	storage.Storage
}

//NewApp returns new app instance
func NewApp(s storage.Storage) *App {
	return &App{s}
}
