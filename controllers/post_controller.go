package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/yona3/go-api-sample/database"
	"github.com/yona3/go-api-sample/ent"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

func (c *PostController) Index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// if exists id
	var id string
	path := strings.Split(r.URL.Path, "/")
	if len(path) == 3 {
		id = path[2]
		log.Println("id = ", id)
	} else if len(path) > 3 {
		log.Println("Too many path segments.")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Too many path segments."))
		return
	}

	switch r.Method {
	case "GET":
		c.getPosts(ctx, w, r)
	case "POST":
		c.createPost(ctx, w, r)
	case "DELETE":
		c.deletePost(ctx, w, r, id)
	default:
		log.Print("Method is not allowed.")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}

// method: GET
func (c *PostController) getPosts(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	db := database.GetClient()

	// get posts
	posts, err := db.Post.Query().All(ctx)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// marshal posts
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
		log.Println("Content-Type must be application/json")
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
		log.Println("'text' and 'user_name' are required.")
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

// method DELETE
func (c *PostController) deletePost(ctx context.Context, w http.ResponseWriter, r *http.Request, id string) {
	// check id
	var postId int
	if id == "" {
		log.Println("id is required.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("'id' is required."))
		return
	} else {
		newId, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("'id' must be integer."))
			return
		}

		postId = newId
	}

	db := database.GetClient()

	// delete post
	err := db.Post.DeleteOneID(postId).Exec(ctx)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	log.Println(fmt.Sprintf("post is deleted. id: %d", postId))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("(id: %d) post is deleted.", postId)))
}
