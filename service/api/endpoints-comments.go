package api

import (
	"encoding/json"
	"net/http"

	"github.com/bibi2001/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) CommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the JSON request body
	var requestBody struct {
		Text string `json:"text"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// Get the text from the JSON request
	text := requestBody.Text
	if !ValidateCommentText(text) {
		http.Error(w, "Invalid caption format", http.StatusBadRequest)
		return
	}

	// Read the photoId from the parameters
	photoId := ps.ByName("photoId")
	// There is no need to validate the identifier
	// we already do that on the database respective file

	// Get the Bearer Token in the header
	token, err := GetBearerToken(r)
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

	// Comment the photo
	comment, err := rt.db.CommentPhoto(authUsername, photoId, text)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't upload comment")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(comment)
}

func (rt *_router) UncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	commentId := ps.ByName("commentId")
	// There is no need to validate the identifier
	// we already do that on the database respective file

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

	isAuthor, err := rt.db.IsAuthor(authUsername, commentId)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't validate comment author")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if !isAuthor {
		http.Error(w, "Authenticated user doesn't match comment's author",
			http.StatusUnauthorized)
		return
	}

	// Uncomment photo
	err = rt.db.UncommentPhoto(commentId)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't uncomment photo")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
