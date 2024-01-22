package handlers

import (
	"cvwo-project/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Comment struct {
	CommentID      int       `json:"comment_id"`
	CommentContent string    `json:"comment_content"`
	UserID         int       `json:"user_id"`
	PostID         int       `json:"post_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// Adds a new comment to a post
func CreateComment(c *gin.Context) {
	var comment Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the referenced user_id and post_id exist in their respective tables
	if !recordExists(c, "users", "user_id", comment.UserID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	if !recordExists(c, "posts", "post_id", comment.PostID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post_id"})
		return
	}

	_, err := database.DB.Exec("INSERT INTO comment (comment_content, user_id, post_id) VALUES ($1, $2, $3)",
		comment.CommentContent, comment.UserID, comment.PostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
 
// Updates an existing comment
func UpdateComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var updatedComment Comment
	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = database.DB.Exec("UPDATE comment SET comment_content = $1 WHERE comment_id = $2",
		updatedComment.CommentContent, commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}


// Deletes a comment by comment_id
func DeleteComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM comment WHERE comment_id = $1", commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Gets all the comments by post_id
func GetCommentsByPostID(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	rows, err := database.DB.Query("SELECT comment_id, comment_content, user_id, post_id, created_at FROM comment WHERE post_id = $1", postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.CommentID, &comment.CommentContent, &comment.UserID, &comment.PostID, &comment.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		comments = append(comments, comment)
	}

	c.JSON(http.StatusOK, comments)
}
