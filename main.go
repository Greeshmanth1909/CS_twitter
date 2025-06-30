package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	fmt.Println("Starting cs_twitter server")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	dbURL := os.Getenv("URL")
	fmt.Println(port, dbURL)
	// open connection to database
	// db, error := sql.Open("postgres", dbURL)
	// if error != nil {
	//     log.Fatal("Error establishing a connection to the database")
	// }

	// dbQueries := database.New(db)
	// apiConf.DB = dbQueries

	// mux := http.NewServeMux()

	// var server http.Server
	// server.Addr = "localhost:" + port
	// server.Handler = mux

	// mux.HandleFunc("GET /v1/healthz", healthHandler)
	// mux.HandleFunc("GET /v1/err", errHandler)
	// mux.HandleFunc("POST /v1/users", usersHandler)
	// mux.Handle("GET /v1/users", authMiddleWare(http.HandlerFunc(getUsersHandler)))
	// mux.Handle("POST /v1/feeds", authMiddleWare(http.HandlerFunc(createFeedHandler)))
	// mux.HandleFunc("GET /v1/feeds", getFeedsHandler)
	// mux.Handle("POST /v1/feed_follows", authMiddleWare(http.HandlerFunc(createFeedFollow)))
	// mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", deleteFeedFollows)
	// mux.Handle("GET /v1/feed_follows", authMiddleWare(http.HandlerFunc(getFeedsByUserId)))

	// server.ListenAndServe()
}
