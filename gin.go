package main

import (
	"cvwo-project/database"
	"cvwo-project/handlers"
	"github.com/gin-gonic/gin"
        "github.com/gin-contrib/cors"
	_ "net/http"
	_ "github.com/lib/pq"
)

func main() {
	database.InitDB()
	router := gin.Default()

        // Apply CORS middleware
        config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Replace with your frontend URL
        config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
        config.AllowHeaders = []string{"Authorization", "Content-Type"}
	router.Use(cors.New(config))

        router.Use(func(c *gin.Context) {
                if c.Request.Method == "OPTIONS" {
                    c.Status(200)
                    c.Abort()
                    return
                }
                c.Next()
            })

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
        router.POST("/likes", handlers.AuthMiddleware(), handlers.CreateLike)
        router.DELETE("deletelike", handlers.AuthMiddleware(), handlers.DeleteLike)
        router.GET("likecount/:post_id", handlers.CountLikesByPostID)

        // Comment Handlers
        router.POST("/createcomment", handlers.AuthMiddleware(), handlers.CreateComment)
        router.PATCH("/updatecomment/:comment_id", handlers.AuthMiddleware(), handlers.UpdateComment)
        router.DELETE("deletecomment/:comment_id", handlers.AuthMiddleware(), handlers.DeleteComment)
        router.GET("/getcommentsbypostid/:post_id", handlers.GetCommentsByPostID)

	router.Run("0.0.0.0:8082")

}
