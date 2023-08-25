package database

func (db *appdbimpl) commentPhoto(username string, photoId int64, text string) (Comment, error) {
	res, err := db.c.Exec(`INSERT INTO comments (id, photoId, username, text) VALUES (NULL, ?, ?, ?)`,
		photoId, username, text)
	if err != nil {
		return nil, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	commentId = uint64(lastInsertID)
	return Comment{commentId, photoId, username, text}, nil
}

func (db *appdbimpl) uncommentPhoto(commentId uint64) error {
	res, err := db.c.Exec(`DELETE FROM comments WHERE commentId=?`, commentId)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If we didn't delete any row, then the comment didn't exist
		return ErrCommentDoesNotExist
	}
	return nil
}