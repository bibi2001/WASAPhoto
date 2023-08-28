package database

func (db *appdbimpl) GetUserStream(username string) ([]Photo, error) {
	var ret []Photo

	// Plain simple SELECT
	rows, err := db.c.Query(`
		SELECT photoId 
		FROM photos AS p, follows AS f
		WHERE p.username = f.followedUser
		AND f.followingUser=?
		ORDER BY p.date DESC`, 
		username)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	// Here we read the resultset and we build the list to be returned
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		p, err := getPhoto(username, id)
		if err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}
