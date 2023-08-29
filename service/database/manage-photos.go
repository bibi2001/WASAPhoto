package database
import "github.com/bibi2001/WASAPhoto/service/globaltime"

func (db *appdbimpl) UploadPhoto(username string, caption string) (Photo, error) {
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

p := Photo{int64(lastInsertID), username, t, caption, 0, 0, false}
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
	err := db.c.Query(`SELECT * FROM photos WHERE photoId=?`, photoId).Scan(
		&p.PhotoId, &p.Username, &p.Date, &p.Caption)
	if err == sql.ErrNoRows {
		return nil, ErrPhotoDoesNotExist
	}	
	if err != nil {
		return nil, err
	}
	
	// Plain simple SELECT to get photo info on comments table
	err := db.c.Query(`SELECT count(*) FROM comments WHERE photoId=?`, photoId).Scan(&p.NComments)
	if err != nil {
		return nil, err
	}
	// Plain simple SELECT to get photo info on likes table
	err := db.c.Query(`SELECT count(*) FROM likes WHERE photoId=?`, photoId).Scan(&p.NLikes)
	if err != nil {
		return nil, err
	}
	// Plain simple SELECT to get photo and user relation info on likes table
	err:= db.c.Query(`SELECT EXISTS (
		SELECT 1 FROM likes WHERE photoId = ? AND username = ?
		)`, photoId, username).Scan(&p.IsLiked)
	if err != nil {
		return nil, err
	}
	
	return p, nil
}

func (db *appdbimpl) ListUserPhotos(username string) ([]Photo, error) {
	var ret []Photo

	// Plain simple SELECT
	rows, err := db.c.Query(`SELECT photoId FROM photos WHERE username=?`, 
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


func (db *appdbimpl) IsPhotoOwner(username string, photoId int64) (bool, error) {
    var isPhotoOwner bool

	// Plain simple SELECT
    err := db.c.QueryRow(`SELECT EXISTS (
        SELECT 1 FROM bans WHERE username=? AND photoId = ?
    )`, username, photoId).Scan(&isPhotoOwner)

    if err != nil {
        return false, err
    }

    return isPhotoOwner, nil
}

