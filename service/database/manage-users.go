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

