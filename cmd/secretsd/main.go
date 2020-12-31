package main

import (
		"flag"
		"fmt"
		"github.com/mkalafior/wp-atrd-task/internal/app/secrets"
		"github.com/mkalafior/wp-atrd-task/internal/mongo"
		"log"
		"net/http"
		"os"
		"os/signal"
		"syscall"
)

func main() {
		addr := flag.String("port", ":3000", "port number")
		flag.Parse()

		// todo hardcoded for demonstration purposes should be read from config or env variables
		db := mongo.NewDb("mongodb://root:root@secrets.mongodb:27017")
		//db := mongo.NewDb("mongodb://root:root@localhost:27017")
		defer mongo.Close()

		sr := mongo.NewSecretRepository(db)
		srv := secrets.NewService(sr)
		errChan := make(chan error)

		go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
				errChan <- fmt.Errorf("%s", <-c)
		}()

		newSecretEndpoint := secrets.MakeNewSecretEndpoint(srv)
		fetchEndpoint := secrets.MakeFetchEndpoint(srv)
		endpoints := secrets.Endpoints{
				NewSecretEndpoint: newSecretEndpoint,
				FetchEndpoint:     fetchEndpoint,
		}

		go func() {
				log.Println("port: ", *addr)
				handler := secrets.NewHttpServer(endpoints)
				http.Handle("/", handler)
				errChan <- http.ListenAndServe(*addr, nil)
		}()

		log.Fatalln(<-errChan)
}
