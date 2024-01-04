DROP TABLE IF EXISTS album;
// Example
CREATE TABLE yourmom (
  id         SERIAL NOT NULL, 
  title      VARCHAR(128) NOT NULL,
  artist     VARCHAR(255) NOT NULL,
  price      DECIMAL(5,2) NOT NULL,
);


// MAIN DATABASE SCHEMA

CREATE TABLE users (
  user_id      SERIAL PRIMARY KEY,
  username     VARCHAR(50) NOT NULL UNIQUE,
);

CREATE TABLE posts (
  post_id       SERIAL PRIMARY KEY,
  title         VARCHAR(255) NOT NULL,
  content       TEXT,
  user_id       INT REFERENCES users(user_id) ON DELETE CASCADE,
  created_at    TIMESTAMP DEFAULT current_timestamp

);

CREATE TABLE comment (
  comment_id    SERIAL PRIMARY KEY,
  content       TEXT NOT NULL,
  user_id       INT REFERENCES users(user_id) ON DELETE CASCADE,
  post_id       INT REFERENCES posts(post_id) ON DELETE CASCADE,
  created_at    TIMESTAMP DEFAULT current_timestamp
  reply_id      //search
  
);


CREATE TABLE post_likes (
  post_id       INT REFERENCES posts(post_id) ON DELETE CASCADE,
  user_id       INT REFERENCES users(user_id) ON DELETE CASCADE,
  
);

CREATE TABLE post_category (
  category_id       SERIAL PRIMARY KEY,
  post_id           INT REFERENCES posts(post_id) ON DELETE CASCADE,
  
);




INSERT INTO yourmom
  (title, artist, price)
VALUES
  ('Blue Train', 'John Coltrane', 56.99),
  ('Giant Steps', 'John Coltrane', 63.99),
  ('Jeru', 'Gerry Mulligan', 17.99),
  ('Sarah Vaughan', 'Sarah Vaughan', 34.98);