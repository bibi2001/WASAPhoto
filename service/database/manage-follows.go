package database

func (db *appdbimpl) followUser(authUser string, followedUser string) error {
	res, err := db.c.Exec(`INSERT INTO follows (authUser, followedUser) VALUES (?, ?)`,
		authUser, followedUser)
	if err != nil {
		return nil, err
	}

	return nil
}

var ErrFollowDoesNotExist = errors.New("The user is not followed!")
func (db *appdbimpl) unfollowUser(authUser string, followedUser string) error {
	res, err := db.c.Exec(`DELETE FROM follows WHERE authUser=? AND followedUser=?`, 
		authUser, followedUser)
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

func (db *appdbimpl) listFollows(authUser string) ([]string, error) {
	var ret []string

	// Plain simple SELECT
	rows, err := db.c.Query(`SELECT followedUser FROM follows WHERE authUser=?`, 
		authUser)
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
