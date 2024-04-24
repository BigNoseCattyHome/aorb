package main

import (
	"log"

	"../aorb/routes"
)

func main() {
	log.Printf("Server started")

	router := routes.NewRouter()

	log.Fatal(router.Run(":8080"))
}
