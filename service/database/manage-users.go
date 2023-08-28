package database

func (db *appdbimpl) CreateUser(username string, name string) error {
	res, err := db.c.Exec(`INSERT INTO users (username, name) VALUES (?, ?)`,
		username, name)
	// Check if username is unique 
	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return errors.New("Username given is not original enough.")
	} else if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) UpdateUsername(oldUsername string, newUsername string) error {
	res, err := db.c.Exec(`UPDATE users SET username=? WHERE username=?`,
		newUsername, oldUsername)
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
	err := db.c.Query(`SELECT * FROM users WHERE username=?`, username).Scan(
		&up.Username, &up.Name)
	if err == sql.ErrNoRows {
		return nil, ErrUserDoesNotExist
	} else if err != nil {
		return nil, err
	}

	// Save number of followers into NFollowers
	followers, err := db.ListFollowers(username)
    if err != nil {
        return nil, err
    }
    up.NFollowers = int64(len(followers))

	// Save number of following into NFollowing
	following, err := db.ListFollowing(username)
	if err != nil {
		return nil, err
	}
	up.NFollowing = int64(len(following))

	// Save if user is followed into isFollowed
	isFollowed, err := db.IsFollowed(authUser, username)
	if err != nil {
		return nil, err
	}
	up.IsFollowed = bool(isFollowed)

	// Save if user is banned into isBanned
	isBanned, err := db.IsBanned(authUser, username)
	if err != nil {
		return nil, err
	}
	up.IsBanned = bool(isBanned)

	// Save user photos into photos
	photos, err := db.ListUserPhotos(username)
	if err != nil {
		return nil, err
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
			SELECT bannedUsername FROM bans 
			WHERE bannedUsername = username AND authUser = ?
		)`, searchQuery, authUser)
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