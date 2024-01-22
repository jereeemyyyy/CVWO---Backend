package handlers

import (
	_"database/sql"
	"cvwo-project/database"
	_ "fmt"
	_"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

type Like struct {
	LikeID int `json:"like_id"`
	PostID int `json:"post_id"`
	UserID int `json:"user_id"`
}

// Adds a new like
func CreateLike(c *gin.Context) {
	var like Like
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the referenced post_id and user_id exist in their respective tables
	if !recordExists(c, "posts", "post_id", like.PostID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post_id"})
		return
	}

	if !recordExists(c, "users", "user_id", like.UserID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	_, err := database.DB.Exec("INSERT INTO likes (post_id, user_id) VALUES ($1, $2)", like.PostID, like.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// Deletes a like by like_id
func DeleteLike(c *gin.Context) {
	likeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid like ID"})
		return
	}

	// Perform the delete operation in the database
	_, err = database.DB.Exec("DELETE FROM likes WHERE like_id = $1", likeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}




// Checks if record exists
func recordExists(c *gin.Context, tableName, columnName string, value int) bool {
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM "+tableName+" WHERE "+columnName+" = $1)", value).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}
	return exists
}

// Counts the number of likes by post_id
func CountLikesByPostID(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Perform the count operation in the database
	var likeCount int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = $1", postID).Scan(&likeCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post_id": postID, "like_count": likeCount})
}