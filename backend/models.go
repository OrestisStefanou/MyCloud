package main

//User json struct
type User struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

//LoginInfo json struct
type LoginInfo struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//DirectoryInfo json struct
type DirectoryInfo struct {
	Path string `json:"path" binding:"required"`
}

//NewDirectoryInfo json struct
type NewDirectoryInfo struct {
	Path    string `json:"path" binding:"required"`
	DirName string `json:"dirName" binding:"required"`
}

//DirectoryEntryInfo json struct
type DirectoryEntryInfo struct {
	Name        string `json:"name" binding:"required"`
	IsDirectory bool   `json:"isDirectory" binding:"required"`
	Icon        string `json:"icon" binding:"required"`
	Link        string `json:"link" binding:"required"`
}

//DirectoryToDeleteInfo json struct
type DirectoryToDeleteInfo struct {
	Name        string `json:"name"`
	IsDirectory bool   `json:"isDirectory"`
	Path        string `json:"path"`
}
