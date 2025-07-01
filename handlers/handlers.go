package handlers

import (
	"context"
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

// The HealthHandler is used to check the health of the server; it sends a 200-OK if the server is running.
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

// The SignupUser handler checks weather a given username is taken and creates a new entry in the users table if it doesn't exist.
func SignupUser(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Username string `json:"username"`
		Password string `json:"password"`
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

// The LoginUser handler sends a jwt if the user exists and the password-hash matches with the one in the database.
func LoginUser(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req body
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&req)

	ctx := context.TODO()
	user, err := apiConf.DB.GetUser(ctx, req.Username)

	// if user doesn't exist
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("username: %v doesn't exist, please sign-up", req.Username)))
		return
	}

	hash := generateHash(req.Password)

	if hash != user.Hash {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid password, try again"))
		return
	}

	// generate jwt with username
	jwt, err := generateJWT(req.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

	resBody, _ := json.Marshal(struct {
		Token string `json:"token"`
	}{Token: jwt})

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resBody)
}
