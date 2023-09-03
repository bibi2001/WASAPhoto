package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Login
	rt.router.POST("/session", rt.wrap(rt.Login))

	// Stream
	rt.router.GET("/home", rt.wrap(rt.GetMyStream))

	// Users
	rt.router.GET("/users/search", rt.wrap(rt.SearchUsers))
	rt.router.GET("/users/:username", rt.wrap(rt.GetUserProfile))
	rt.router.PUT("/users/:username/username", rt.wrap(rt.SetMyUsername))

	// Follows
	rt.router.GET("/users/:username/following", rt.wrap(rt.ListFollowing))
	rt.router.GET("/users/:username/followers", rt.wrap(rt.ListFollowers))
	rt.router.PUT("/users/:username/followers/:authUser", rt.wrap(rt.FollowUnfollow))
	rt.router.DELETE("/users/:username/followers/:authUser", rt.wrap(rt.FollowUnfollow))

	// Bans
	rt.router.GET("/users/:username/bans", rt.wrap(rt.ListBans))
	rt.router.PUT("/users/:username/bans/:bannedUser", rt.wrap(rt.BanUnban))
	rt.router.DELETE("/users/:username/bans/:bannedUser", rt.wrap(rt.BanUnban))

	// Photos
	rt.router.POST("/photo", rt.wrap(rt.UploadPhoto))
	rt.router.DELETE("/photo/:photoId", rt.wrap(rt.DeletePhoto))
	rt.router.GET("/photo/:photoId", rt.wrap(rt.GetPhoto))

	// Comments
	rt.router.GET("/photo/:photoId/comments", rt.wrap(rt.ListComments))
	rt.router.POST("/photo/:photoId/comments", rt.wrap(rt.CommentPhoto))
	rt.router.DELETE("/photo/:photoId/comments/:commentId", rt.wrap(rt.UncommentPhoto))

	// Likes
	rt.router.GET("/photo/:photoId/likes", rt.wrap(rt.ListLikes))
	rt.router.PUT("/photo/:photoId/likes/:username", rt.wrap(rt.LikeUnlike))
	rt.router.DELETE("/photo/:photoId/likes/:username", rt.wrap(rt.LikeUnlike))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
