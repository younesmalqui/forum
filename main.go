package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/api"
	"forum/config"
	"forum/handlers"
	"forum/utils"
)

func main() {
	if err := utils.InitServices(); err != nil {
		log.Fatal(err)
	}
	if err := utils.InitTables(); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/static/", handlers.ServeStatic)
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/filter", handlers.FilterHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogutHandler)
	http.HandleFunc("/post/{id}", handlers.PostHandler)
	http.HandleFunc("/api/register", api.RegisterApi)
	http.HandleFunc("/api/post", api.PostApi)
	http.HandleFunc("/api/react", api.ReactToPostHandler)
	http.HandleFunc("/api/add/comment", api.AddComment)
	http.HandleFunc("/api/like/comment", api.HandleLikeComment)
	fmt.Printf("Server running on http://localhost%v", config.ADDRS)
	err := http.ListenAndServe(config.ADDRS, nil)
	if err != nil {
		log.Fatal(err)
	}
}
