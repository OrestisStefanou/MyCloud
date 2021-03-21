package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const clientsBaseDir = "/home/orestis/MyCloud"

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
	}

	router.Run()
}

func signup(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err == nil {
		err = createUser(newUser)
		if err == nil {
			c.JSON(http.StatusAccepted, gin.H{"status": "created"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A user with this email already exists"})
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are necessary"})
	}
}

func signin(c *gin.Context) {
	var userLoginInfo LoginInfo
	if err := c.ShouldBindJSON(&userLoginInfo); err == nil {
		user, err := getUser(userLoginInfo.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong"})
		} else {
			if user.Email != "" { //If user exists
				if user.Password == userLoginInfo.Password {
					session := sessions.Default(c)
					session.Set("userEmail", user.Email)
					session.Save()
					c.JSON(http.StatusAccepted, gin.H{"status": "Login successful"})
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong username or password"})
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong username or password"})
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are necessary"})
	}
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"status": "Log out successfull"})
}

func isAuthenticated(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("userEmail")
	if email != nil {
		c.JSON(http.StatusAccepted, gin.H{"status": "Authenticated", "email": email.(string)})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authenticated"})
	}
}

func listClientDir(c *gin.Context) {
	session := sessions.Default(c)
	var dirInfo DirectoryInfo
	var userEmail string
	email := session.Get("userEmail")
	if email != nil {
		userEmail = email.(string)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authenticated"})
	}

	if err := c.ShouldBindJSON(&dirInfo); err == nil {
		userPath := filepath.Join(clientsBaseDir, userEmail, dirInfo.Path)
		dir, err := ioutil.ReadDir(userPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong!"})
		}
		fmt.Println("Listing subdir/parent")
		for _, entry := range dir {
			fmt.Println(" ", entry.Name(), entry.IsDir())
		}

		c.JSON(http.StatusAccepted, gin.H{"status": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path data not given"})
	}
}

func uploadFile(c *gin.Context) {
	session := sessions.Default(c)
	var userEmail string
	email := session.Get("userEmail")
	if email != nil {
		userEmail = email.(string)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authenticated"})
	}
	// single file
	file, _ := c.FormFile("file")
	path := c.PostForm("path")
	log.Println(file.Filename, path)

	// Upload the file to specific dst.
	//dst := "/home/orestis/Desktop/git.pdf"
	dst := filepath.Join(clientsBaseDir, userEmail, path, file.Filename)
	fmt.Println("FILEPATH TO CREATE ", dst)
	c.SaveUploadedFile(file, dst)

	c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("'%s' uploaded!", file.Filename)})
}

func createDirectory(c *gin.Context) {
	session := sessions.Default(c)
	var newDirInfo NewDirectoryInfo
	var userEmail string
	email := session.Get("userEmail")
	if email != nil {
		userEmail = email.(string)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authenticated"})
	}

	if err := c.ShouldBindJSON(&newDirInfo); err == nil {
		newDir := filepath.Join(clientsBaseDir, userEmail, newDirInfo.Path, newDirInfo.DirName)
		fmt.Println("CREATING DIRECTORY ", newDir)
		c.JSON(http.StatusOK, gin.H{"status": "Directory created"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not all neccesarry data given"})
	}
}
