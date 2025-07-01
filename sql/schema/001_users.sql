-- +goose Up
CREATE TABLE USERS (
    username varchar(255) UNIQUE NOT NULL,
    hash varchar(255) NOT NULL
);

CREATE TABLE POSTS (
    post_id uuid DEFAULT gen_random_uuid(),
    post TEXT NOT NULL,
    username varchar(255),
    PRIMARY KEY (post_id),
    FOREIGN KEY (username) REFERENCES USERS(username)
);

CREATE TABLE COMMENTS (
    comment_id uuid DEFAULT gen_random_uuid(),
    comment TEXT NOT NULL,
    post_id uuid,
    username varchar(255),
    PRIMARY KEY (comment_id),
    FOREIGN KEY (post_id) REFERENCES POSTS(post_id),
    FOREIGN KEY (username) REFERENCES USERS(username)
);

-- +goose Down
DROP TABLE COMMENTS;
DROP TABLE POSTS;
DROP TABLE USERS;