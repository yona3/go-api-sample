package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/yona3/go-api-sample/database"
	"github.com/yona3/go-api-sample/ent"
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
	case "POST":
		c.createPost(ctx, w, r)
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

// method: POST
func (c *PostController) createPost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	db := database.GetClient()

	// check content-type
	h := r.Header.Get("Content-Type")
	if h != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Content-Type must be application/json"))
		return
	}

	// read body
	b := r.Body
	defer b.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, b)

	var post *ent.Post
	if err := json.Unmarshal(buf.Bytes(), &post); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// check required fields
	if post.Text == "" || post.UserName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("'text' and 'user_name' are required."))
		return
	}

	// create post
	res, err := db.Post.Create().SetText(post.Text).SetUserName(post.UserName).Save(ctx)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	log.Println(fmt.Sprintf("new post: %+v", res))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("(id: %d) post is created.", res.ID)))
}
