package handlers

import (
	"database/sql"
	"cvwo-project/database"
	_ "fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

type Post struct {
	PostID    int       `db:"post_id" json:"post_id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	UserID    int       `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Create a new post
func CreatePost(c *gin.Context) {
	var post Post
	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the current timestamp
	post.CreatedAt = time.Now()

	// Insert the post into the database
	err := database.DB.QueryRow(
        `INSERT INTO posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING post_id`,
        post.Title, post.Content, post.UserID,
    ).Scan(&post.PostID)

    if err != nil {
        log.Println(err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
        return
    }

    c.JSON(http.StatusCreated, post)
}

// Get post by ID
func GetPostByID(c *gin.Context) {
	postIDstr := c.Param("id")
	postID, err := strconv.Atoi(postIDstr)
	log.Printf("Attempting to fetch post with ID: %d", postID)


	var post Post
	err = database.DB.QueryRow("SELECT * FROM posts WHERE post_id = $1", postID).Scan(
		&post.PostID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}

		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// Get all posts
func GetAllPosts(c *gin.Context) {
	var posts []Post

	rows, err := database.DB.Query("SELECT * FROM posts")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.PostID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// Update a post

func UpdatePost(c *gin.Context) {
    postID := c.Param("id")

    // Check if the post exists
    var existingPost Post
    err := database.DB.QueryRow("SELECT * FROM posts WHERE post_id = $1", postID).
        Scan(&existingPost.PostID, &existingPost.Title, &existingPost.Content, &existingPost.UserID, &existingPost.CreatedAt)

    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
        return
    }

    var updatedPost Post
    if err := c.ShouldBindJSON(&updatedPost); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Update only the fields that are provided
    if updatedPost.Title != "" {
        existingPost.Title = updatedPost.Title
    }
    if updatedPost.Content != "" {
        existingPost.Content = updatedPost.Content
    }
    if updatedPost.UserID != 0 {
        existingPost.UserID = updatedPost.UserID
    }

    // Update the post in the database
    _, err = database.DB.Exec(`
        UPDATE posts
        SET title = $1, content = $2, user_id = $3
        WHERE post_id = $4
    `, existingPost.Title, existingPost.Content, existingPost.UserID, existingPost.PostID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
        return
    }

    c.JSON(http.StatusOK, existingPost)
}



// Delete a post

func DeletePost(c *gin.Context) {
	postID := c.Param("id")

	// Check if the post exists
	var existingPost Post
	err := database.DB.QueryRow("SELECT * FROM posts WHERE post_id = $1", postID).Scan(
		&existingPost.PostID, &existingPost.Title, &existingPost.Content, &existingPost.UserID, &existingPost.CreatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
		return
	}

	// Delete the post from the database
	_, err = database.DB.Exec("DELETE FROM posts WHERE post_id = $1", postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}