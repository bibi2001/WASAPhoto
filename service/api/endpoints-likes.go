package api

import (
	"net/http"

	"github.com/bibi2001/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) LikeUnlike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Read the photoId from the parameters
	photoId := ps.ByName("photoId")
	authUser := ps.ByName("username")

	// Check the validaty of the photoId parameter
	photoExists, err := rt.db.PhotoExists(photoId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error searching for the photo")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !photoExists {
		http.Error(w, "Invalid photoID", http.StatusBadRequest)
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

	// Check if we should like or unlike the photo
	if r.Method == http.MethodPut {
		// Like the photo
		err = rt.db.LikePhoto(authUser, photoId)
		if err != nil {
			ctx.Logger.WithError(err).Error("can't like photo")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else if r.Method == http.MethodDelete {
		// Unlike the photo
		err = rt.db.UnlikePhoto(authUser, photoId)
		if err != nil {
			ctx.Logger.WithError(err).Error("can't unlike photo")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

}
