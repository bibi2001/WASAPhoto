package database

import "errors"

func (db *appdbimpl) LikePhoto(userId string, photoId int64) error {
	_, err := db.c.Exec(`INSERT INTO likes (photoId, userId) VALUES (?, ?)`,
		photoId, userId)
	if err != nil {
		return err
	}

	return nil
}

var ErrLikeDoesNotExist = errors.New("the like does not exist")

func (db *appdbimpl) UnlikePhoto(userId string, photoId int64) error {
	res, err := db.c.Exec(`DELETE FROM likes WHERE photoId=? AND userId=?`, photoId, userId)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If we didn't delete any row, then the like didn't exist
		return ErrLikeDoesNotExist
	}
	return nil
}

func (db *appdbimpl) ListLikes(photoId int64) ([]string, error) {
	var ret []string

	// Plain simple SELECT
	rows, err := db.c.Query(`SELECT username FROM likes WHERE photoId=?`, photoId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	// Here we read the resultset and we build the list to be returned
	for rows.Next() {
		var l string
		err = rows.Scan(&l)
		if err != nil {
			return nil, err
		}

		ret = append(ret, l)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}
