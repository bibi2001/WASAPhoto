package api

import (
	"encoding/json"
	"net/http"

	"github.com/bibi2001/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) UploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the JSON request body
	var requestBody struct {
		Image   string `json:"image"`
		Caption string `json:"caption"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// Get the image data from the JSON request
	imageBytes := []byte(requestBody.Image)

	// Get the caption from the JSON request
	caption := requestBody.Caption

	if !ValidateCaption(caption) {
		http.Error(w, "Invalid caption format", http.StatusBadRequest)
		return
	}

	// Get the Bearer Token in the header
	token, err := GetBearerToken(r)
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

	// Upload the photo
	photo, err := rt.db.UploadPhoto(authUser, caption, imageBytes)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't upload photo")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(photo)
}
