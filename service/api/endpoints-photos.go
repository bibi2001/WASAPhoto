package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/bibi2001/WASAPhoto/service/api/reqcontext"
	"github.com/bibi2001/WASAPhoto/service/utils"
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

	if !utils.ValidateCaption(caption) {
		http.Error(w, "Invalid caption format", http.StatusBadRequest)
		return
	}

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

	// Upload the photo to the database
	photo, err := rt.db.UploadPhoto(authUser, caption)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't upload photo to database")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Define the file path where we'll save the image
	filePath := "./storage/" + strconv.FormatInt(photo.PhotoId, 10) + ".jpg"

	// Create or open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't create file for photo")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Write the image bytes to the file
	_, err = file.Write(imageBytes)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't write image data to file")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(photo)
}

func (rt *_router) DeletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Read the photoId from the parameters
	photoId := ps.ByName("photoId")
	// There is no need to validate the identifier as we already do that on the database respective file

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

	// Check if authenticated user is the owner of the photo
	isPhotoOwner, err := rt.db.IsPhotoOwner(authUser, utils.ToInt64(photoId))
	if err != nil {
		ctx.Logger.WithError(err).Error("could not get photo owner")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !isPhotoOwner {
		// User isn't owner of the photo, so they can't delete it
		http.Error(w, "Not the owner of the photo", http.StatusUnauthorized)
		return
	}

	err = rt.db.DeletePhoto(utils.ToInt64(photoId))
	if err != nil {
		ctx.Logger.WithError(err).Error("could not delete photo")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) GetPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Read the photoId from the parameters
	photoId := ps.ByName("photoId")
	// There is no need to validate the identifier as we already do that on the database respective file

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

	// Get the photo from the database
	photo, err := rt.db.GetPhoto(authUser, utils.ToInt64(photoId))
	if err != nil {
		ctx.Logger.WithError(err).Error("could not get user photo")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return photo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(photo)

}
