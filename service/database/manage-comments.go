package database

import (
	"errors"

	"github.com/bibi2001/WASAPhoto/service/api/utils"
)

func (db *appdbimpl) CommentPhoto(username string, photoId int64, text string) (Comment, error) {
	res, err := db.c.Exec(`INSERT INTO comments (id, photoId, username, text) VALUES (NULL, ?, ?, ?)`,
		photoId, username, text)
	if err != nil {
		return nil, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	commentId := uint64(lastInsertID)
	return Comment{commentId, photoId, username, text}, nil
}

var ErrCommentDoesNotExist = errors.New("The comment does not exist!")

func (db *appdbimpl) UncommentPhoto(commentId uint64) error {
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

func (db *appdbimpl) ListComments(photoId int64) ([]Comment, error) {
	var ret []SimpleComment

	// Plain simple SELECT
	rows, err := db.c.Query(`SELECT commentId, photoId, username, text FROM comments WHERE photoId=?`,
		photoId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	// Here we read the resultset and we build the list to be returned
	for rows.Next() {
		var c Comment
		err = rows.Scan(&c.CommentId, &c.PhotoId, &c.Username, &c.Text)
		if err != nil {
			return nil, err
		}

		ret = append(ret, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (db *appdbimpl) IsAuthor(username string, commentId int64) (bool, error) {
	var isAuthor bool

	// Plain simple SELECT
	err := db.c.QueryRow(`SELECT EXISTS (
        SELECT 1 FROM photos WHERE username=? AND photoId = ?
    )`, username, commentId).Scan(&isAuthor)

	if err != nil {
		return false, err
	}

	return isAuthor, nil
}
