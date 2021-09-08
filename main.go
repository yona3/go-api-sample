package main

import (
	"context"
	"log"
	"net/http"

	"github.com/yona3/go-api-sample/controllers"
	"github.com/yona3/go-api-sample/database"
)

const PORT = "8080"

func main() {
	database.Init()
	log.Println("Connected to database")

	defer database.CloseClient()

	ctx := context.Background()
	postController := controllers.NewPostController()

	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) { postController.Index(ctx, w, r) })
	log.Println("server running on port " + PORT)
	http.ListenAndServe(":"+PORT, nil)
}
