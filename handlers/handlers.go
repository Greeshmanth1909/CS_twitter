package handlers

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
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

func SignupUser(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Username string `json:username`
		Password string `json:password`
	}
	var req body
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&req)

	ctx := context.TODO()
	_, err := apiConf.DB.GetUser(ctx, req.Username)

	if err == nil {
		// A user with existing username was found!
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(fmt.Sprintf("username: %v already taken", req.Username)))
		return
	}

	passwordHash := generateHash(req.Password)

	// add to the database
	var addUserParams database.AddUserParams
	addUserParams.Username = req.Username
	addUserParams.Hash = passwordHash

	user, err := apiConf.DB.AddUser(ctx, addUserParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

	res, _ := json.Marshal(user)
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Username string `json:username`
		Password string `json:password`
	}
	var req body
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&req)

	var User database.User
	User.Username = req.Username

	hash := sha256.New()
	hash.Write([]byte(req.Password))
	hashed := hash.Sum(nil)

	User.Hash = string(hashed)

	w.WriteHeader(200)
}
