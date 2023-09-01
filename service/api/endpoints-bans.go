package api

import (
	"encoding/json"
	"net/http"

	"github.com/bibi2001/WASAPhoto/service/api/reqcontext"
	"github.com/bibi2001/WASAPhoto/service/utils"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) BanUnban(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Read the usernames from the parameters
	bannedUser := ps.ByName("bannedUser")
	authUser := ps.ByName("username")

	if bannedUser == authUser {
		// Here we validate the usernames, and we
		// discovered that they are not valid.
		// A user can't ban or unban themselves
		http.Error(w, "Invalid Usernames. A user can't ban or unban themselves",
			http.StatusBadRequest)
		return
	}

	usernameExists, err := rt.db.UserExists(bannedUser)
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
	token, err := utils.GetBearerToken(r)
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

	// Check if we should ban or unban the user
	if r.Method == http.MethodPut {
		// Ban user
		err = rt.db.BanUser(authUser, bannedUser)
		if err != nil {
			ctx.Logger.WithError(err).Error("can't ban user")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else if r.Method == http.MethodDelete {
		// Unban user
		err = rt.db.UnbanUser(authUser, bannedUser)
		if err != nil {
			ctx.Logger.WithError(err).Error("can't unban user")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

}

func (rt *_router) ListBans(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Read the username from the parameters
	username := ps.ByName("username")

	// Get the Bearer Token in the header
	token, err := utils.GetBearerToken(r)
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
	// Check if the authenticated user is the profile owner
	// Only the profile owner can check bans list
	if username != authUsername {
		http.Error(w, "Authenticated userID doesn't match profile owner",
			http.StatusUnauthorized)
		return
	}

	// Get the bans list from the database
	bansList, err := rt.db.ListBans(username)
	if err != nil {
		ctx.Logger.WithError(err).Error("could not get bans list")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return bans list
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(bansList)
}
