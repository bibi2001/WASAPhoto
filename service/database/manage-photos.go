package database
import "github.com/bibi2001/WASAPhoto/service/globaltime"

func (db *appdbimpl) uploadPhoto(username string, caption string) (Photo, error) {
	t := globaltime.Now()
	res, err := db.c.Exec(`INSERT INTO photos (photoId, username, date, caption) 
		VALUES (NULL, ?, ?, ?)`, username, t, caption)
if err != nil {
	return nil, err
}

lastInsertID, err := res.LastInsertId()
if err != nil {
	return nil, err
}

p := Photo{uint64(lastInsertID), username, t, caption, 0, 0, false}
return p, nil
}

var ErrPhotoDoesNotExist = errors.New("The photo does not exist!")
func (db *appdbimpl) deletePhoto(photoId int64) error {
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

func (db *appdbimpl) getPhoto(username string, photoId int64) (Photo, error) {

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
}
