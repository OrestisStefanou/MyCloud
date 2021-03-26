package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	dbConnect("orestis", "Ore$tis1997", "MyCloud")
	router := gin.Default()

	//Session middleware
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	//Cors middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	v1 := router.Group("/v1/MyCloud")
	{
		v1.POST("/signup", signup)
		v1.POST("/signin", signin)
		v1.GET("/logout", logout)
		v1.GET("/isAuthenticated", isAuthenticated)
		v1.POST("/upload", uploadFile)
		v1.POST("/listDir", listClientDir)
		v1.POST("/createDir", createDirectory)
		v1.POST("/delete", deleteEntry)
		v1.StaticFS("/static", http.Dir("./static"))
		v1.StaticFS("/files", http.Dir(clientsBaseDir))
	}

	router.Run()
}
