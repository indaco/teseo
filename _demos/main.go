package main

import (
	"log"
	"net/http"
	"os"

	"github.com/indaco/teseo/_demos/handlers"
)

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./_demos/statics"))
	mux.Handle("GET /statics/", http.StripPrefix("/statics/", fileServer))

	mux.HandleFunc("GET /", handlers.HandleHome)
	mux.HandleFunc("GET /company", handlers.HandleCompany)
	mux.HandleFunc("GET /about", handlers.HandleAbout)
	mux.HandleFunc("GET /blog", handlers.HandleBlog)
	mux.HandleFunc("GET /blog/posts/{id}", handlers.HandlePosts)

	port := ":3300"
	log.Printf("Listening on %s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Printf("failed to start server: %v", err)
		os.Exit(1)
	}
}
