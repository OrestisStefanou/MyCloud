package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
		urlPath := getPathFromURLParam(dirInfo.Path) //Create the url path from url parameter
		userPath := filepath.Join(clientsBaseDir, userEmail, urlPath)
		dir, err := ioutil.ReadDir(userPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong!"})
		}
		var dirEntries = make([]DirectoryEntryInfo, len(dir))
		for i, entry := range dir {
			//fmt.Println(" ", entry.Name(), entry.IsDir())
			dirEntries[i].Icon = getFileIconLink(entry.Name(), entry.IsDir())
			dirEntries[i].Link = getFileLink(filepath.Join(userEmail, urlPath, entry.Name()))
			dirEntries[i].Name = entry.Name()
			dirEntries[i].IsDirectory = entry.IsDir()
		}

		c.JSON(http.StatusAccepted, gin.H{"dirEntries": dirEntries})
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
	urlPath := c.PostForm("path")
	path := getPathFromURLParam(urlPath)
	log.Println(file.Filename, path)

	// Upload the file to specific dst.
	dst := filepath.Join(clientsBaseDir, userEmail, path, file.Filename)
	c.SaveUploadedFile(file, dst)

	var newFileInfoResponse DirectoryEntryInfo
	newFileInfoResponse.Name = file.Filename
	newFileInfoResponse.IsDirectory = false
	newFileInfoResponse.Icon = getFileIconLink(file.Filename, false)
	newFileInfoResponse.Link = getFileLink(filepath.Join(userEmail, path, file.Filename))
	c.JSON(http.StatusOK, gin.H{"newFile": newFileInfoResponse})
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
		newDir := filepath.Join(clientsBaseDir, userEmail, getPathFromURLParam(newDirInfo.Path), newDirInfo.DirName)
		err = os.Mkdir(newDir, 0755)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong during creating the folder"})
		} else {
			var newDirInfoResponse DirectoryEntryInfo
			newDirInfoResponse.Name = newDirInfo.DirName
			newDirInfoResponse.IsDirectory = true
			newDirInfoResponse.Icon = getFileIconLink(newDirInfo.DirName, true)
			newDirInfoResponse.Link = "TestLink"
			c.JSON(http.StatusOK, gin.H{"newDirectory": newDirInfoResponse})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Directory name not given"})
	}
}

func deleteEntry(c *gin.Context) {
	session := sessions.Default(c)
	var userEmail string
	email := session.Get("userEmail")
	if email != nil {
		userEmail = email.(string)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authenticated"})
	}

	var dirEntry DirectoryToDeleteInfo
	if err := c.ShouldBindJSON(&dirEntry); err == nil {
		pathToDelete := filepath.Join(clientsBaseDir, userEmail, dirEntry.Path, dirEntry.Name)
		if dirEntry.IsDirectory {
			err := os.RemoveAll(pathToDelete)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Internal Error"})
			}
			c.JSON(http.StatusOK, gin.H{"status": "Deleted"})
		} else {
			err := os.Remove(pathToDelete)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Internal Error"})
			}
			c.JSON(http.StatusOK, gin.H{"status": "Deleted"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error during parsing the object"})
	}
}
