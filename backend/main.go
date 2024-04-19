package main

import (
	"log"

	"github.com/BigNoseCattyHome/aorb/routes"
)

func main() {
	log.Printf("Server started")

	router := routes.NewRouter()

	log.Fatal(router.Run(":8080"))
}
