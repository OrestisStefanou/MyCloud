package main

type User struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type LoginInfo struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type DirectoryInfo struct {
	Path string `json:"path" binding:"required"`
}

type NewDirectoryInfo struct {
	Path    string `json:"path" binding:"required"`
	DirName string `json:"dirName" binding:"required"`
}

type DirectoryEntryInfo struct {
	Name        string `json:"name" binding:"required"`
	IsDirectory bool   `json:"isDirectory" binding:"required"`
	Icon        string `json:"icon" binding:"required"`
	Link        string `json:"link" binding:"required"`
}

type DirectoryToDeleteInfo struct {
	Name        string `json:"name"`
	IsDirectory bool   `json:"isDirectory"`
	Path        string `json:"path"`
}
