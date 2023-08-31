package database

import (
	"database/sql"
	"errors"

	"github.com/bibi2001/WASAPhoto/service/utils"
	"github.com/bibi2001/WASAPhoto/service/globaltime"
)

func (db *appdbimpl) UploadPhoto(username string, caption string, image []byte) (Photo, error) {
	t := globaltime.Now()

	res, err := db.c.Exec(`INSERT INTO photos (photoId, image, username, date, caption) 
		VALUES (NULL, ?, ?, ?, ?)`, image, username, t, caption)
	if err != nil {
		return Photo{}, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return Photo{}, err
	}

	p := Photo{PhotoId: int64(lastInsertID), Image: image, Username: username,
		Date: t, Caption: caption, NComments: 0, NLikes: 0, IsLiked: false}
	return p, nil
}

var ErrPhotoDoesNotExist = errors.New("The photo does not exist!")

func (db *appdbimpl) DeletePhoto(photoId int64) error {
	res, err := db.c.Exec(`DELETE FROM photos WHERE photoId=?`, photoId)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If we didn't delete any row, then the like didn't exist
		return ErrPhotoDoesNotExist
	}
	return nil
}

func (db *appdbimpl) GetPhoto(username string, photoId int64) (Photo, error) {
	var p Photo

	// Plain simple SELECT to get photo info on photos table
	err := db.c.QueryRow(`SELECT photoId, image, username, date, caption 
		FROM photos WHERE photoId=?`, photoId).Scan(
		&p.PhotoId, &p.Image, &p.Username, &p.Date, &p.Caption)
	if err == sql.ErrNoRows {
		return Photo{}, ErrPhotoDoesNotExist
	}
	if err != nil {
		return Photo{}, err
	}

	// Plain simple SELECT to get photo info on comments table
	err = db.c.QueryRow(`SELECT count(*) FROM comments WHERE photoId=?`, photoId).Scan(&p.NComments)
	if err != nil {
		return Photo{}, err
	}
	// Plain simple SELECT to get photo info on likes table
	err = db.c.QueryRow(`SELECT count(*) FROM likes WHERE photoId=?`, photoId).Scan(&p.NLikes)
	if err != nil {
		return Photo{}, err
	}
	// Plain simple SELECT to get photo and user relation info on likes table
	err = db.c.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM likes WHERE photoId = ? AND username = ?
		)`, photoId, username).Scan(&p.IsLiked)
	if err != nil {
		return Photo{}, err
	}

	return p, nil
}

func (db *appdbimpl) ListUserPhotos(username string) ([]Photo, error) {
	var ret []Photo

	// Plain simple SELECT
	rows, err := db.c.Query(`SELECT photoId FROM photos
		WHERE username = ? ORDER BY date DESC`, username)
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
		p, err := db.GetPhoto(username, id)
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

func (db *appdbimpl) IsPhotoOwner(username string, photoId int64) (bool, error) {
	var isPhotoOwner bool

	// Plain simple SELECT
	err := db.c.QueryRow(`SELECT EXISTS (
        SELECT 1 FROM photos WHERE username=? AND photoId = ?
    )`, username, photoId).Scan(&isPhotoOwner)

	if err != nil {
		return false, err
	}

	return isPhotoOwner, nil
}
