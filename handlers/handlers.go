package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/Greeshmanth1909/CS_twitter/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

var apiConf apiConfig

// Initialize a persistant connection to local database with init function
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("URL")
	// open connection to database
	db, error := sql.Open("postgres", dbURL)
	if error != nil {
		fmt.Println(error)
		log.Fatal("Error establishing a connection to the database")
	}

	fmt.Println(db)

	dbQueries := database.New(db)
	apiConf.DB = dbQueries
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	dbRes, err := apiConf.DB.ListUsers(ctx)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
	}

	fmt.Println(dbRes)
}
