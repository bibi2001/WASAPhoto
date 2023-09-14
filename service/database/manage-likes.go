package database

import "errors"

func (db *appdbimpl) LikePhoto(username string, photoId int64) error {
	_, err := db.c.Exec(`INSERT INTO likes (photoId, username) VALUES (?, ?)`,
		photoId, username)
	if err != nil {
		return err
	}

	return nil
}

var ErrLikeDoesNotExist = errors.New("the like does not exist")

func (db *appdbimpl) UnlikePhoto(username string, photoId int64) error {
	res, err := db.c.Exec(`DELETE FROM likes WHERE photoId=? AND username=?`, photoId, username)
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
