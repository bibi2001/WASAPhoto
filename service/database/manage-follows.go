package database

import "errors"

func (db *appdbimpl) FollowUser(followingUser string, followedUser string) error {
	_, err := db.c.Exec(`INSERT INTO follows (followingUser, followedUser) VALUES (?, ?)`,
		followingUser, followedUser)
	if err != nil {
		return err
	}

	return nil
}

var ErrFollowDoesNotExist = errors.New("the user is not followed")

func (db *appdbimpl) UnfollowUser(followingUser string, followedUser string) error {
	res, err := db.c.Exec(`DELETE FROM follows WHERE followingUser=? AND followedUser=?`,
		followingUser, followedUser)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If we didn't delete any row, then the follow didn't exist
		return ErrFollowDoesNotExist
	}
	return nil
}

func (db *appdbimpl) ListFollowRelationship(query string, username string) ([]string, error) {
	var ret []string

	// Plain simple SELECT
	rows, err := db.c.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	// Here we read the resultset and we build the list to be returned
	for rows.Next() {
		var f string
		err = rows.Scan(&f)
		if err != nil {
			return nil, err
		}

		ret = append(ret, f)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (db *appdbimpl) ListFollowing(followingUser string) ([]string, error) {
	return db.ListFollowRelationship(`SELECT followedUser FROM follows WHERE followingUser=?`, followingUser)
}

func (db *appdbimpl) ListFollowers(followedUser string) ([]string, error) {
	return db.ListFollowRelationship(`SELECT followingUser FROM follows WHERE followedUser=?`, followedUser)
}

func (db *appdbimpl) IsFollowed(followingUser string, followedUser string) (bool, error) {
	var isFollowed bool

	// Plain simple SELECT
	err := db.c.QueryRow(`SELECT EXISTS (
        SELECT 1 FROM follows WHERE followingUser=? AND followedUser = ?
    )`, followingUser, followedUser).Scan(&isFollowed)

	if err != nil {
		return false, err
	}

	return isFollowed, nil
}
