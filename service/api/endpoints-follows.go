package api

import (
	"net/http"

	"github.com/bibi2001/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) FollowUnfollow(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Read the usernames from the parameters
	userToFollow := ps.ByName("username")
	authUser := ps.ByName("authUser")

	if userToFollow == authUser {
		// Here we validate the usernames, and we
		// discovered that they are not valid.
		// A user can't follow or unfollow themselves
		http.Error(w, "Invalid Usernames. A user can't follow or unfollow themselves",
			http.StatusBadRequest)
		return
	}

	usernameExists, err := rt.db.UserExists(userToFollow)
	if err != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		ctx.Logger.WithError(err).Error("can't validate username")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !usernameExists {
		// Here we validate the username, and we
		// discovered that it is not valid.
		http.Error(w, "Invalid Username", http.StatusBadRequest)
		return
	}

	// Get the Bearer Token in the header
	token, err := GetBearerToken(r)
	if err != nil {
		http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
		return
	}
	// Get the username corresponding to the Token
	authUsername, err := rt.db.GetUsername(token)
	if err != nil {
		ctx.Logger.WithError(err).Error("authenticated username can not be found")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if authUser != authUsername {
		http.Error(w, "Authenticated userID doesn't match authenticated username",
			http.StatusUnauthorized)
		return
	}

	// Check if we should follow or unfollow the user
	if r.Method == http.MethodPut {
		// Follow user
		err = rt.db.FollowUser(authUser, userToFollow)
		if err != nil {
			ctx.Logger.WithError(err).Error("can't follow user")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else if r.Method == http.MethodDelete {
		// Unfollow user
		err = rt.db.UnfollowUser(authUser, userToFollow)
		if err != nil {
			ctx.Logger.WithError(err).Error("can't unfollow user")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

}
