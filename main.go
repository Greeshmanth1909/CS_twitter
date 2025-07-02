package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/Greeshmanth1909/CS_twitter/handlers"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Starting cs_twitter server")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	mux := http.NewServeMux()

	var server http.Server
	server.Addr = "localhost:" + port
	server.Handler = mux

	fileServer := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fileServer)

	mux.HandleFunc("GET /v1/health", handlers.HealthHandler)
	mux.HandleFunc("GET /v1/posts", handlers.GetPosts)
	mux.HandleFunc("POST /v1/login", handlers.LoginUser)
	mux.HandleFunc("POST /v1/signup", handlers.SignupUser)
	mux.Handle("POST /v1/create-post", handlers.AuthMiddleWare(http.HandlerFunc(handlers.CreatePost)))
	mux.Handle("POST /v1/create-comment", handlers.AuthMiddleWare(http.HandlerFunc(handlers.CreateComment)))
	server.ListenAndServe()
}
