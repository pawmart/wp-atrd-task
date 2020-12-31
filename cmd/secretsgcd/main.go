package main

import (
		"fmt"
		"github.com/mkalafior/wp-atrd-task/internal/mongo"
		"log"
		"os"
		"os/signal"
		"syscall"
		"time"
)

func main() {
		errChan := make(chan error)

		go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
				errChan <- fmt.Errorf("%s", <-c)
		}()

		db := mongo.NewDb("mongodb://root:root@secrets.mongodb:27017")
		defer mongo.Close()

		repo := mongo.NewSecretGCRepository(db)
		go func() {
				for {
						err := repo.GC()
						if err != nil {
								errChan <- err
						}
						log.Println("GC cycle finished")
						time.Sleep(60 * time.Second)
				}
		}()

		log.Fatalln(<-errChan)
}
