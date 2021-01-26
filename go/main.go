package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	addRoutes(r)
	address := ":" + getEnv("PORT", "8080")
	go func() {
		err := http.ListenAndServe(address, r)
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("Listening on", address)
	select {}
}
