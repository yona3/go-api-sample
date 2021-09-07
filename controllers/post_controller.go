package controllers

import (
	"encoding/json"
	"net/http"
)

type Post struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
	UserName  string `json:"userName"`
}

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

func (c *PostController) Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		c.getPosts(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}

// method: GET
func (c *PostController) getPosts(w http.ResponseWriter, _ *http.Request) {
	Posts := []Post{
		{ID: 1, Text: "Hello!", CreatedAt: "2021-01-01 20:00", UserName: "John"},
		{ID: 2, Text: "Hey.", CreatedAt: "2021-01-01 20:28", UserName: "Tom"},
		{ID: 3, Text: "I'm yona.", CreatedAt: "2021-01-02 12:02", UserName: "yona"},
	}

	jsonBytes, err := json.Marshal(Posts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
