package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/Greeshmanth1909/CS_twitter/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

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

// The GetPosts handler returns all the posts, corresponding comments and associated data like username etc.
func GetPosts(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	feed, err := apiConf.DB.GetFeed(ctx)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
	}

	type ResponseStruct struct {
		Username string
		Post     string
		Post_id  uuid.UUID
		Comments [][]string
	}
	var res []ResponseStruct

	for i := range feed {
		new := ResponseStruct{}
		new.Username = feed[i].UsernamePost
		new.Post = feed[i].Post
		new.Post_id = feed[i].PostID

		// process aggregated strings; get rid of {} and ""
		c := string(feed[i].Comments.([]uint8))
		d := string(feed[i].CommenterUsernames.([]uint8))
		c = strings.TrimPrefix(c, "{")
		c = strings.TrimSuffix(c, "}")
		stringSlice := strings.Split(c, ",")

		// for j := range stringSlice {
		// 	if len(stringSlice[j]) >= 2 {
		// 		stringSlice[j] = stringSlice[j][1 : len(stringSlice[j])-1]
		// 	}
		// }

		d = strings.TrimPrefix(d, "{")
		d = strings.TrimSuffix(d, "}")

		comments_u := Zip(stringSlice, strings.Split(d, ",")) // zip comments and corresponding usernames
		new.Comments = comments_u
		res = append(res, new)
	}

	resp, _ := json.Marshal(res)
	w.Write(resp)
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

	_, err = apiConf.DB.AddUser(ctx, addUserParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

	w.WriteHeader(http.StatusCreated)
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

	// generate jwt with username included in the claims
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

// The CreatePost handler creates a post by uploading it to the database along with the username from middleware. This is a protected route i.e. needs jwt from the login handler.
func CreatePost(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(jwtClaims).(jwt.MapClaims)
	username := claims["username"].(string)

	type body struct {
		Post string `json:"post"`
	}
	var req body
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&req)

	var createUserPostParams database.CreateUserPostParams
	createUserPostParams.Username = username
	createUserPostParams.Post = req.Post

	ctx := context.TODO()
	res, err := apiConf.DB.CreateUserPost(ctx, createUserPostParams)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

	responseJson, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(responseJson)
}

// The CreateComments handler creates a comment on a post, it needs the post's post_id. It is a also a protected route.
func CreateComment(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(jwtClaims).(jwt.MapClaims)
	username := claims["username"].(string)

	type body struct {
		Comment string `json:"comment"`
		PostID  string `json:"post_id"`
	}
	var req body
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&req)

	var createUserCommentParams database.CreateUserCommentParams
	createUserCommentParams.Username = username
	createUserCommentParams.Comment = req.Comment
	createUserCommentParams.PostID, _ = uuid.Parse(req.PostID)

	ctx := context.TODO()
	comment, err := apiConf.DB.CreateUserComment(ctx, createUserCommentParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}
	resBody, _ := json.Marshal(comment)

	w.WriteHeader(http.StatusAccepted)
	w.Write(resBody)
}
