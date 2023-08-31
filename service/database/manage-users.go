package database

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/bibi2001/WASAPhoto/service/utils"
)

func (db *appdbimpl) CreateUser(username string) error {
	_, err := db.c.Exec(`INSERT INTO users (userId, username) VALUES (NULL, ?)`, username)
	// Check if username is unique
	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return errors.New("Username given is not original enough.")
	} else if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) UpdateUsername(oldUsername string, newUsername string) error {
	_, err := db.c.Exec(`UPDATE users SET username=? WHERE username=?`, newUsername, oldUsername)
	// Check if username is unique
	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return errors.New("Username given is not original enough.")
	} else if err != nil {
		return err
	}

	return nil
}

var ErrUserDoesNotExist = errors.New("The user does not exist!")

func (db *appdbimpl) GetUserProfile(username string, authUser string) (UserProfile, error) {
	var up UserProfile

	// Plain simple SELECT to get user info on users table
	err := db.c.QueryRow(`SELECT * FROM users WHERE username=?`, username).Scan(
		&up.Username, &up.UserId)
	if err == sql.ErrNoRows {
		return UserProfile{}, ErrUserDoesNotExist
	} else if err != nil {
		return UserProfile{}, err
	}

	// Save number of followers into NFollowers
	followers, err := db.ListFollowers(username)
	if err != nil {
		return UserProfile{}, err
	}
	up.NFollowers = int64(len(followers))

	// Save number of following into NFollowing
	following, err := db.ListFollowing(username)
	if err != nil {
		return UserProfile{}, err
	}
	up.NFollowing = int64(len(following))

	// Save if user is followed into isFollowed
	isFollowed, err := db.IsFollowed(authUser, username)
	if err != nil {
		return UserProfile{}, err
	}
	up.IsFollowed = isFollowed

	// Save if user is banned into isBanned
	isBanned, err := db.IsBanned(authUser, username)
	if err != nil {
		return UserProfile{}, err
	}
	up.IsBanned = isBanned

	// Save user photos into photos
	photos, err := db.ListUserPhotos(username)
	if err != nil {
		return UserProfile{}, err
	}
	up.Photos = photos
	up.NPosts = int64(len(photos))

	return up, nil
}

func (db *appdbimpl) UserSearch(searchQuery string, authUser string) ([]string, error) {
	var ret []string

	// Simple SELECT query to get matches that aren't banned by the authenticated user
	rows, err := db.c.Query(`
		SELECT username FROM users  
		WHERE username LIKE '%' || ? || '%'
		AND username NOT IN (
			SELECT bannedUser FROM bans 
			WHERE bannedUser = username 
			AND authUser=? ) 
		AND username NOT IN (
			SELECT authUser FROM bans 
			WHERE authUser = username 
			AND bannedUser=?
		)`, searchQuery, authUser, authUser)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	// Here we read the resultset and we build the list to be returned
	for rows.Next() {
		var u string
		err = rows.Scan(&u)
		if err != nil {
			return nil, err
		}

		ret = append(ret, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (db *appdbimpl) UserExists(username string) (bool, error) {
	var exists bool

	err := db.c.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM users WHERE username=?
	)`, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (db *appdbimpl) GetUserId(username string) (int64, error) {
	userExists, err := db.UserExists(username)
	if err != nil {
		return 0, err
	}

	if !userExists {
		// Attempt to create the user
		if err := db.CreateUser(username); err != nil {
			return 0, err
		}
	}

	var userId int64
	err = db.c.QueryRow(`SELECT userId FROM users WHERE username=?`, username).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (db *appdbimpl) GetUsername(userId int64) (string, error) {
	var username string

	err := db.c.QueryRow(`SELECT username FROM users WHERE userId=?`, userId).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}
