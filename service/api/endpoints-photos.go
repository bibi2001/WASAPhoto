package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bibi2001/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) UploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the multipart form with a file size limit of 32Mb
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error opening file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the image into a byte slice
	imageBytes, err := ioutil.ReadAll(file)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed reading the image file")
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}

	// Get the caption from the form
	caption := r.FormValue("caption")
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
		ctx.Logger.WithError(err).Error("authenticated username can not be found")
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
