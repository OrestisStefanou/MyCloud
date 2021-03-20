package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func dbConnect(dbUser, dbPass, dbName string) {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPass, dbName))
	checkErr(err)

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

func getUser(email string) (User, error) {
	user := User{}
	rows, err := db.Query("SELECT * FROM Users WHERE Email=?", email)
	if err != nil {
		return user, err
	}
	for rows.Next() {
		err = rows.Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password)
		if err != nil {
			return user, err
		}
	}
	return user, nil
}

func createUser(user User) error {
	stmt, err := db.Prepare("INSERT Users SET FirstName=?,LastName=?,Email=?,Password=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}
	//Create a directory for the user
	err = os.MkdirAll(filepath.Join(clientsBaseDir, user.Email, "home"), 0755)
	if err != nil {
		return err
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
