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

func (db *appdbimpl) getUserProfile(username string) error {
	var p Photo 
	// Plain simple SELECT to get photo info on photos table
	err := db.c.Query(`SELECT * FROM photos WHERE photoId=?`, photoId).Scan(
		&p.photoId, &p.username, &p.date, &p.caption)
	if err == sql.ErrNoRows {
		return nil, ErrPhotoDoesNotExist
	}	
	if err != nil {
		return nil, err
	}
	
	// Plain simple SELECT to get photo info on comments table
	err := db.c.Query(`SELECT count(*) FROM comments WHERE photoId=?`, photoId).Scan(&p.nComments)
	if err != nil {
		return nil, err
	}
	// Plain simple SELECT to get photo info on likes table
	err := db.c.Query(`SELECT count(*) FROM likes WHERE photoId=?`, photoId).Scan(&p.nlikes)
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
