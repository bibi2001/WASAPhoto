/*
The structs.go contains the useful structs that are used on the rest of the service directory,
more specifically in the database and api subdirectories.
*/

package utils

import "time"

type Comment struct {
	CommentId int64  `json:"commentId"`
	PhotoId   int64  `json:"photoId"`
	Username  string `json:"username"`
	Text      string `json:"text"`
}

type Photo struct {
	PhotoId   int64     `json:"photoId"`
	Username  string    `json:"username"`
	Image     []byte    `json:"image"`
	Date      time.Time `json:"date"`
	Caption   string    `json:"caption"`
	NComments int64     `json:"nComments"`
	NLikes    int64     `json:"nLikes"`
	IsLiked   bool      `json:"isLiked"`
}

type UserProfile struct {
	Username   string  `json:"username"`
	UserId     string  `json:"userId"`
	NPosts     int64   `json:"nPosts"`
	NFollowers int64   `json:"nFollowers"`
	NFollowing int64   `json:"nFollowing"`
	IsFollowed bool    `json:"isFollowed"`
	IsBanned   bool    `json:"isbanned"`
	Photos     []Photo `json:"photos"`
}

type LoginResponse struct {
	Identifier string `json:"identifier"`
}
