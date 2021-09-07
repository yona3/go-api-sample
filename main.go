package main

import (
	"log"
	"net/http"

	"github.com/yona3/go-api-sample/controllers"
)

const PORT = "8080"

func main() {
	postController := controllers.NewPostController()

	http.HandleFunc("/todos", postController.Index)
	log.Println("server running on port " + PORT)
	http.ListenAndServe(":"+PORT, nil)
}
