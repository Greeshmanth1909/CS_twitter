# CS_twitter
CS_twitter is a simple application that lets users Post text based posts and comment on them.

## Build from source
### Requirements
- Go `version 1.24`
- Postgres
- `sqlc` to generate database queries
- `goose` to manage database migrations
- `Docker`

### Build instructions
1. Clone the repo with `git clone https://github.com/Greeshmanth1909/CS_twitter.git`
2. Install dependency with `go mod tidy`
3. Run the postgres docker container with `docker-compose up`
4. Run database migrations with `make migrate-up` and `make migrate-down`
5. Run `sqlc generate` to generate database queries for go. This creates a `internal/database` directory.
6. Run `go build .` or `make build` in the root directory to build the binary
7. Set the port of the server in the `.env` file the default is 8080.
8. Run with `./main` or `make run`

### Usage
1. visit `localhost:8080/` to access the frontend.
2. Please login/signup to post/comment.
3. Rich text support has been added with markdown syntax.
```
For bold, place sentence between a pair of two asterisks (**) eg. one **bold** word in sentence.
For italic, place sentence between a pair of asterisks (*) eg. one *italic* word in sentence.
For Hyperlinks, place the text to be displayed between square brackets [] followed by the http/https link between parenthesis (). eg. [just google it](https://www.google.com/)
```

## Routes

### Signup User
```http
POST /v1/signup
Content-Type: application/json

{
    "username": "your_username",
    "password": "your_password"
}
```

### Login
```http
POST /v1/login
Content-Type: application/json

{
    "username": "your_username",
    "password": "your_password"
}
```
Response includes a JWT token that should be used in the Authorization header for protected routes:
```
Authorization: Bearer <your_jwt_token>
```
### Health
```http
GET /v1/health
```
Responds with `200-OK`

### GetPosts
```http
GET /v1/GetPosts
```
Returns all posts and corresponding comments
```http
[{"Username":"user","Post":"Hello there! ","Comments":[["",""]]},{"Username":"User1","Post":"second post","Comments":[["This is a comment","User"],["This is another comment","User"],["This comment by diff user","user1"]]}]
```

### Create Post (Protected route)
```http
POST /v1/create-posts
Content-Type: application/json

{
  "post": "second post"
}
```
Responds with post details and 202 status code
```
{"PostID":"3579872a-ddd1-4736-b531-2d3fe43b44cb","Post":"second post","Username":"User"}
```

### Create Comment (Protected route)
```http
POST /v1/create-comment
Content-Type: application/json

{
  "comment": "comment",
  "post_id": "uuid of post"
}
```
Responds with post details and 202 status code
```
{"CommentID":"e15487e1-95fd-4cfb-b810-ce5c8b8d9968","Comment":"comment","PostID":"uuid of post","Username":"user1"}

```

## Architecture
### Database
Since all the data that needs to be stored is text based and is structured into posts and comments, a relational database like Postgres was chosen. The application uses Postgres running in a docker container with docker compose.
`sqlc` generates go code to query and update the database.
`goose` is used to manage all database migrations.

### Programming language
Go was choosen for its simplicity and robust standard library.



