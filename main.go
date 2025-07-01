package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/Greeshmanth1909/CS_twitter/handlers"
	"github.com/Greeshmanth1909/CS_twitter/internal/database"

	_ "github.com/lib/pq"
)

func main() {

	fmt.Println("Starting cs_twitter server")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	dbURL := os.Getenv("URL")
	fmt.Printf("Connecting to db with port %v and url %v\n", port, dbURL)

	// open connection to database
	db, error := sql.Open("postgres", dbURL)
	if error != nil {
		fmt.Println(error)
		log.Fatal("Error establishing a connection to the database")
	}

	fmt.Println(db)

	dbQueries := database.New(db)
	apiConf.DB = dbQueries

	mux := http.NewServeMux()

	var server http.Server
	server.Addr = "localhost:" + port
	server.Handler = mux

	mux.HandleFunc("GET /v1/health", handlers.HealthHandler)
	mux.HandleFunc("GET /v1/posts", handlers.GetPosts)
	// mux.HandleFunc("POST /v1/login", usersHandler)
	// mux.HandleFunc("POST /v1/signup", usersHandler)
	// mux.Handle("GET /v1/users", authMiddleWare(http.HandlerFunc(getUsersHandler)))
	// mux.Handle("POST /v1/feeds", authMiddleWare(http.HandlerFunc(createFeedHandler)))
	// mux.HandleFunc("GET /v1/feeds", getFeedsHandler)
	// mux.Handle("POST /v1/feed_follows", authMiddleWare(http.HandlerFunc(createFeedFollow)))
	// mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", deleteFeedFollows)
	// mux.Handle("GET /v1/feed_follows", authMiddleWare(http.HandlerFunc(getFeedsByUserId)))

	server.ListenAndServe()
}

type apiConfig struct {
	DB *database.Queries
}

var apiConf apiConfig
