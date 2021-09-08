package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/yona3/go-api-sample/database"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

func (c *PostController) Index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		c.getPosts(ctx, w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}

// method: GET
func (c *PostController) getPosts(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	db := database.GetClient()

	posts, err := db.Post.Query().All(ctx)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	jsonBytes, err := json.Marshal(posts)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	log.Println(string(jsonBytes))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
