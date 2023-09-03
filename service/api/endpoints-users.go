package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bibi2001/WASAPhoto/service/api/reqcontext"
	"github.com/bibi2001/WASAPhoto/service/utils"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Read the username from the parameters
	username := ps.ByName("username")

	// Check if the requested username exists
	usernameExists, err := rt.db.UserExists(username)
	if err != nil {
		// Log the error and return a 500 Internal Server Error
		ctx.Logger.WithError(err).Error("can't validate username")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !usernameExists {
		// Return a 400 Bad Request if the username is invalid
		http.Error(w, "Invalid Username", http.StatusBadRequest)
		return
	}

	// Get the Bearer Token from the request header
	token, err := utils.GetBearerToken(r)
	if err != nil {
		http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
		return
	}

	// Get the username corresponding to the Token
	authUser, err := rt.db.GetUsername(token)
	if err != nil {
		ctx.Logger.WithError(err).Error("authenticated username cannot be found")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get the user profile
	userProfile, err := rt.db.GetUserProfile(username, authUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("could not get user profile")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return the user profile as JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(userProfile)
}

func (rt *_router) SetMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the JSON request body
	var requestBody struct {
		Username string `json:"username"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// Get the new username from the JSON request
	newUsername := requestBody.Username
	if !utils.ValidateUsername(newUsername) {
		http.Error(w, "Invalid username format", http.StatusBadRequest)
		return
	}
	// Read the old username from the parameters
	oldUsername := ps.ByName("username")

	// Check if the new username is different from the old one
	if oldUsername == newUsername {
		http.Error(w, "New username cannot match old username", http.StatusUnauthorized)
		return
	}

	// Get the Bearer Token from the request header
	token, err := utils.GetBearerToken(r)
	if err != nil {
		http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
		return
	}
	// Get the username corresponding to the Token
	authUsername, err := rt.db.GetUsername(token)
	if err != nil {
		ctx.Logger.WithError(err).Error("authenticated username cannot be found")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if the user trying to change the username is authorized to do so
	if authUsername != oldUsername {
		http.Error(w, "Only the owner of the profile can edit the username", http.StatusUnauthorized)
		return
	}

	// Update the username
	err = rt.db.UpdateUsername(oldUsername, newUsername)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't update username")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rt *_router) SearchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Read the photoId from the parameters
	query := r.URL.Query().Get("q")

	// Print the query variable for debugging purposes
	fmt.Println("Search query:", query)
	// Get the Bearer Token in the header
	token, err := utils.GetBearerToken(r)
	if err != nil {
		http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
		return
	}
	// Get the username corresponding to the Token
	authUser, err := rt.db.GetUsername(token)
	if err != nil {
		ctx.Logger.WithError(err).Error("authenticated username cannot be found")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get the search result from the database
	searchResult, err := rt.db.UserSearch(query, authUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("could not search result")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return search result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(searchResult)

}
