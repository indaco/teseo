package handlers

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/indaco/teseo/_demos/pages/posts"
)

func HandlePosts(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// Map of post IDs to corresponding render functions
	postHandlers := map[string]func() templ.Component{
		"first-post":  posts.FirstPostPage,
		"second-post": posts.SecondPostPage,
	}

	if handler, ok := postHandlers[id]; ok {
		if err := handler().Render(r.Context(), w); err != nil {
			log.Printf("Error rendering page: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		http.NotFound(w, r)
	}
}
