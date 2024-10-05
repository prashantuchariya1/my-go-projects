package main

import (
	"fmt"
	"log"
	"mongodbv2/router"
	"net/http"
)

func main() {
	fmt.Println("Welcome to MongoDB V2")
	r := router.Router()
	fmt.Println("Server is getting started...")
	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Println("Listening at port 4000...")

}
