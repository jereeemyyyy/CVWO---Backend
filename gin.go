package main

import (
	"cvwo-project/database"
	"cvwo-project/handlers"
	"github.com/gin-gonic/gin"
	_ "net/http"
	_ "strconv"
	_ "github.com/lib/pq"
)

func main() {
	database.InitDB()
	router := gin.Default()

        // Post handlers
	router.POST("/posts", handlers.AuthMiddleware(), handlers.CreatePost)
        router.GET("/posts/:id", handlers.GetPostByID)
        router.GET("/posts", handlers.GetAllPosts)
        router.DELETE("/posts/:id", handlers.AuthMiddleware(), handlers.DeletePost)
        router.PATCH("/posts/:id", handlers.AuthMiddleware(), handlers.UpdatePost)
        
        // User Handlers
        router.POST("/register", handlers.Register)
        router.POST("/login", handlers.Login)
        router.GET("/userprofile/:id", handlers.AuthMiddleware(), handlers.GetUserProfile)
        router.DELETE("/deleteuser/:id", handlers.AuthMiddleware(), handlers.DeleteUser)

        // Like Handlers
        router.POST("/likes", handlers.CreateLike)
        router.DELETE("deletelike/:id", handlers.AuthMiddleware(), handlers.DeleteLike)
        router.GET("likecount/:post_id", handlers.AuthMiddleware(), handlers.CountLikesByPostID)

        // Comment Handlers
        router.POST("/createcomment", handlers.AuthMiddleware(), handlers.CreateComment)
        router.PATCH("/updatecomment/:comment_id", handlers.AuthMiddleware(), handlers.UpdateComment)
        router.DELETE("deletecomment/:comment_id", handlers.AuthMiddleware(), handlers.DeleteComment)
        router.GET("/getcommentsbypostid/:post_id", handlers.GetCommentsByPostID)

	router.Run("0.0.0.0:80")

}
