/* DATABASE SCHEMA */

CREATE TABLE users (
  user_id      SERIAL PRIMARY KEY,
  username     VARCHAR(50) NOT NULL UNIQUE,
  password_hash     VARCHAR(255) NOT NULL
);

CREATE TABLE posts (
  post_id       SERIAL PRIMARY KEY,
  title         VARCHAR(255) NOT NULL,
  content       TEXT,
  user_id       INT REFERENCES users(user_id) ON DELETE CASCADE,
  created_at    TIMESTAMP DEFAULT current_timestamp

);

CREATE TABLE tags (
    tag_id SERIAL PRIMARY KEY,
    tag_name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE post_tags (
    tag_id INT REFERENCES tags(tag_id) ON DELETE CASCADE,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE likes (
  like_id       SERIAL PRIMARY KEY,
  post_id       INT REFERENCES posts(post_id) ON DELETE CASCADE,
  user_id       INT REFERENCES users(user_id) ON DELETE CASCADE
  
);


CREATE TABLE comment (
  comment_id        SERIAL PRIMARY KEY,
  comment_content   TEXT NOT NULL,
  user_id           INT REFERENCES users(user_id) ON DELETE CASCADE,
  post_id           INT REFERENCES posts(post_id) ON DELETE CASCADE,
  created_at        TIMESTAMP DEFAULT current_timestamp
  
);



