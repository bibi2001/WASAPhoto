package database

func (db *appdbimpl) createUser(username string, name string) error {
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

func (db *appdbimpl) updateUsername(oldUsername string, newUsername string) error {
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
type UserProfile struct {
	Username   string  `json:"username"`
	Name       string  `json:"name"`
	NPosts     int64   `json:"nPosts"`
	NFollowers int64   `json:"nFollowers"`
	NFollowing int64   `json:"nFollowing"`
	IsFollowed bool    `json:"isFollowed"`
	IsBanned   bool    `json:"isbanned"`
	Photos     []Photo `json:"photos"`
}

var ErrUserDoesNotExist = errors.New("The user does not exist!")
func (db *appdbimpl) getUserProfile(username string) error {
	var up UserProfile
	// Plain simple SELECT to get user info on users table
	err := db.c.Query(`SELECT * FROM users WHERE username=?`, username).Scan(
		&up.Username, &up.User)
	if err == sql.ErrNoRows {
		return nil, ErrUserDoesNotExist
	}	
	if err != nil {
		return nil, err
	}
	
	// Plain simple SELECT to get user info on photos table
	err := db.c.Query(`SELECT count(*) FROM photos WHERE username=?`, username).Scan(&up.NPosts)
	if err != nil {
		return nil, err
	}
	// Plain simple SELECT to get user info on comments table
	err := db.c.Query(`SELECT count(*) FROM comments WHERE username=?`, username).Scan(&up.NComments)
	if err != nil {
		return nil, err
	}
	// Plain simple SELECT to get photo info on likes table
	err := db.c.Query(`SELECT count(*) FROM likes WHERE username=?`, username).Scan(&up.NLikes)
	if err != nil {
		return nil, err
	}
	// Plain simple SELECT to get photo and user relation info on likes table
	err:= db.c.Query(`SELECT EXISTS (
		SELECT 1 FROM likes WHERE photoId = 0 AND username = 'jonhDoe'
		)`, photoId, username).Scan(&p.isLiked)
	if err != nil {
		return nil, err
	}
	
	return p, nil

	return nil
}
