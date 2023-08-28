/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	uploadPhoto(username string, caption string) (Photo, error)
	deletePhoto(photoId int64) error
	getPhoto(username string, photoId int64) (Photo, error)
	listUserPhotos(username string) ([]Photo, error)

	commentPhoto(username string, photoId int64, text string) (Comment, error)
	uncommentPhoto(commentId uint64) error 
	listComments(photoId int64) ([]Comment, error)

	likePhoto(username string, photoId int64) error
	unlikePhoto(username string, photoId int64) error
	listLikes(photoId int64) (username []string, error)

	banUser(authUsername string, bannedUser string) error
	unbanUser(authUsername string, bannedUser string) error
	listBans(authUsername string) ([]string, error)

	followUser(followingUser string, followedUser string) error
	unfollowUser(followingUser string, followedUser string) error
	listFollowRelationship(query string, username string) ([]string, error)
	listFollowing(followingUser string) ([]string, error)
	listFollowers(followedUser string) ([]string, error)

	
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Enable foreign keys
	stmt = "PRAGMA foreign_keys = ON"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}
	
	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string

	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE "users" (
			"username"	TEXT NOT NULL,
			"name"	TEXT,
			PRIMARY KEY("username")
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}
	
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='photos';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE "photos" (
			"photoId"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"username"	TEXT NOT NULL,
			"date"	TEXT NOT NULL,
			"caption"	TEXT,
			FOREIGN KEY("username") REFERENCES "users"("username")
				ON DELETE CASCADE
        		ON UPDATE CASCADE
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='comments';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE "comments" (
			"commentId"	INTEGER NOT NULL AUTOINCREMENT,
			"photoId"	INTEGER NOT NULL,
			"username"	TEXT NOT NULL,
			"text"	TEXT NOT NULL,
			PRIMARY KEY("commentId"),
			FOREIGN KEY("photoId") REFERENCES "photos"("photoId")
				ON DELETE CASCADE
        		ON UPDATE CASCADE,
			FOREIGN KEY("username") REFERENCES "users"("username")
				ON DELETE CASCADE
				ON UPDATE CASCADE
			
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='likes';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE "likes" (
			"username"	TEXT NOT NULL,
			"photoId"	INTEGER NOT NULL,
			PRIMARY KEY("username","photoId"),
			FOREIGN KEY("username") REFERENCES "users"("username")
				ON DELETE CASCADE
				ON UPDATE CASCADE,
			FOREIGN KEY("photoId") REFERENCES "photos"("photoId")
				ON DELETE CASCADE
				ON UPDATE CASCADE
			
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='follows';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE "follows" (
			"followingUser"	TEXT NOT NULL,
			"followedUser"	TEXT NOT NULL,
			PRIMARY KEY("followedUser","followingUser"),
			FOREIGN KEY("followingUser") REFERENCES "users"("username")
				ON DELETE CASCADE
				ON UPDATE CASCADE,
			FOREIGN KEY("followedUser") REFERENCES "users"("username")
				ON DELETE CASCADE
				ON UPDATE CASCADE
			
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='bans';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE "bans" (
			"authUser"	TEXT NOT NULL,
			"bannedUser"	TEXT NOT NULL,
			PRIMARY KEY("authUser","bannedUser"),
			FOREIGN KEY("authUser") REFERENCES "users"("username")
				ON DELETE CASCADE
				ON UPDATE CASCADE,
			FOREIGN KEY("bannedUser") REFERENCES "users"("username")
				ON DELETE CASCADE
				ON UPDATE CASCADE
			
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
