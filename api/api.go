package main

import (
	"data"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

// CRUD Route Handlers
func createPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
	decoder := json.NewDecoder(r.Body)
	var newPost data.Post
	if err := decoder.Decode(&newPost); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	data.DB.Create(&newPost)
	res, err := json.Marshal(newPost)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(res)
}

func updatePostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setCors(w)
	type body struct {
		Author  string
		Message string
	}
	var updates body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updates); err != nil {
		http.Error(w, err.Error(), 400)
	}

	var updatedPost data.Post
	data.DB.Where("ID = ?", ps.ByName("postId")).First(&updatedPost)
	updatedPost.Author = updates.Author
	updatedPost.Message = updates.Message
	data.DB.Save(&updatedPost)
	res, err := json.Marshal(updatedPost)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.Write(res)
}

func deletePostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setCors(w)
	var deletedPost data.Post
	data.DB.Where("ID = ?", ps.ByName("postId")).Delete(&deletedPost)
	res, err := json.Marshal(deletedPost)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	w.Write(res)
}

func showPostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setCors(w)
	var post data.Post
	data.DB.Where("ID = ?", ps.ByName("postId")).First(&post)
	res, err := json.Marshal(post)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(res)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
	var posts []data.Post
	data.DB.Find(&posts)
	res, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(res)
}

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
	fmt.Fprintf(w, "This is the RESTful api")
}

// used for COR preflight checks
func corsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
}

// util
func getFrontendURL() string {
	if os.Getenv("APP_ENV") == "production" {
		return "http://localhost:80" // change this to production domain
	} else {
		return "http://localhost:3000"
	}
}

func setCors(w http.ResponseWriter) {
	frontendURL := getFrontendURL()
	w.Header().Set("Access-Control-Allow-Origin", frontendURL)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
	defer data.DB.Close()

	// add router and routes
	router := httprouter.New()
	router.GET("/", indexHandler)
	router.POST("/api/posts", createPostHandler)
	router.GET("/api/posts/:postId", showPostHandler)
	router.DELETE("/api/posts/:postId", deletePostHandler)
	router.PUT("/api/posts/:postId", updatePostHandler)
	router.GET("/api/posts", indexPostHandler)
	router.OPTIONS("/*any", corsHandler)

	env := os.Getenv("APP_ENV")
	fmt.Println(env)
	if env == "production" {
		log.Println("Running api server in production mode")
	} else {
		log.Println("Running api server in dev mode")
	}
	http.ListenAndServe(":4000", router)
}
