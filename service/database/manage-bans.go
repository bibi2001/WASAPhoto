package database

func (db *appdbimpl) BanUser(authUser string, bannedUser string) error {
	res, err := db.c.Exec(`INSERT INTO bans (authUser, bannedUser) VALUES (?, ?)`,
		authUser, bannedUser)
	if err != nil {
		return err
	}
	err := db.UnfollowUser(authUser, bannedUser)
	if err!= nil && err!=db.ErrFollowDoesNotExist {
		return err
	}

	return nil
}

var ErrBanDoesNotExist = errors.New("The user is not banned!")
func (db *appdbimpl) UnbanUser(authUser string, bannedUser string) error {
	res, err := db.c.Exec(`DELETE FROM bans WHERE authUser=? AND bannedUser=?`, 
		authUser, bannedUser)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If we didn't delete any row, then the ban didn't exist
		return ErrBanDoesNotExist
	}
	return nil
}

func (db *appdbimpl) ListBans(authUser string) ([]string, error) {
	var ret []string

	// Plain simple SELECT
	rows, err := db.c.Query(`SELECT bannedUser FROM bans WHERE authUser=?`, 
		authUser)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	// Here we read the resultset and we build the list to be returned
	for rows.Next() {
		var b string
		err = rows.Scan(&b)
		if err != nil {
			return nil, err
		}

		ret = append(ret, b)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (db *appdbimpl) IsBanned(authUser string, bannedUser string) (bool, error) {
    var isBanned bool

	// Plain simple SELECT
    err := db.c.QueryRow(`SELECT EXISTS (
        SELECT 1 FROM bans WHERE authUser=? AND bannedUser = ?
    )`, authUser, bannedUser).Scan(&isBanned)

    if err != nil {
        return false, err
    }

    return isBanned, nil
}
