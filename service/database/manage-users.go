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
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("username given is not original enough")
		}
		return err
	}

	return nil
}

func (db *appdbimpl) UpdateUsername(oldUsername string, newUsername string) error {
	result, err := db.c.Exec(`UPDATE users SET username=? WHERE username=?`, newUsername, oldUsername)

	if err != nil {
		return err // Handle the general database error.
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("username not found")
	}

	return nil
}

var ErrUserDoesNotExist = errors.New("the user does not exist")

func (db *appdbimpl) GetUserProfile(username string, authUser string) (utils.UserProfile, error) {
	var up utils.UserProfile

	// Plain simple SELECT to get user info on users table
	err := db.c.QueryRow(`SELECT username, userId FROM users WHERE username=?`, username).Scan(
		&up.Username, &up.UserId)
	if errors.Is(err, sql.ErrNoRows) {
		return utils.UserProfile{}, ErrUserDoesNotExist
	} else if err != nil {
		return utils.UserProfile{}, err
	}

	// Save number of followers into NFollowers
	followers, err := db.ListFollowers(username)
	if err != nil {
		return utils.UserProfile{}, err
	}
	up.NFollowers = int64(len(followers))

	// Save number of following into NFollowing
	following, err := db.ListFollowing(username)
	if err != nil {
		return utils.UserProfile{}, err
	}
	up.NFollowing = int64(len(following))

	// Save if user is followed into isFollowed
	isFollowed, err := db.IsFollowed(authUser, username)
	if err != nil {
		return utils.UserProfile{}, err
	}
	up.IsFollowed = isFollowed

	// Save if user is banned into isBanned
	isBanned, err := db.IsBanned(authUser, username)
	if err != nil {
		return utils.UserProfile{}, err
	}
	up.IsBanned = isBanned

	// Save user photos into photos
	photos, err := db.ListUserPhotos(username)
	if err != nil {
		return utils.UserProfile{}, err
	}
	up.Photos = photos
	up.NPosts = int64(len(photos))

	return up, nil
}

func (db *appdbimpl) UserSearch(searchQuery string, authUser string) ([]string, error) {
	var ret []string

	searchQuery = "%" + searchQuery + "%"
	// Simple SELECT query to get matches that aren't banned by the authenticated user
	rows, err := db.c.Query(`
		SELECT username FROM users  
		WHERE username LIKE ?
		AND username NOT IN (
			SELECT bannedUser FROM bans 
			WHERE bannedUser = username 
			AND authUser=? ) 
		AND username NOT IN (
			SELECT authUser FROM bans 
			WHERE authUser = username 
			AND bannedUser=?
		AND username != "?"
		)`, searchQuery, authUser, authUser, authUser)
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

func (db *appdbimpl) GetUserId(username string) (string, error) {
	userExists, err := db.UserExists(username)
	if err != nil {
		return "", err
	}

	if !userExists {
		// Attempt to create the user
		if err := db.CreateUser(username); err != nil {
			return "", err
		}
	}

	var userId string
	err = db.c.QueryRow(`SELECT userId FROM users WHERE username=?`, username).Scan(&userId)
	if err != nil {
		return "", err
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
