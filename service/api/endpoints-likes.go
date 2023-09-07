package api

import (
	"encoding/json"
	"net/http"

	"github.com/bibi2001/WASAPhoto/service/api/reqcontext"
	"github.com/bibi2001/WASAPhoto/service/utils"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) LikeUnlike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Read the photoId from the parameters
	photoId := ps.ByName("photoId")
	// There is no need to validate the identifier
	// we already do that on the database respective file
	authUser := ps.ByName("username")

	// Check if we should like or unlike the photo
	if r.Method == http.MethodPut {
		// Like the photo
		err := rt.db.LikePhoto(authUser, utils.ToInt64(photoId))
		if err != nil {
			ctx.Logger.WithError(err).Error("can't like photo")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else if r.Method == http.MethodDelete {
		// Unlike the photo
		err := rt.db.UnlikePhoto(authUser, utils.ToInt64(photoId))
		if err != nil {
			ctx.Logger.WithError(err).Error("can't unlike photo")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

}

func (rt *_router) ListLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Read the photoId from the parameters
	photoId := ps.ByName("photoId")

	// Get the likes list from the database
	likeList, err := rt.db.ListLikes(utils.ToInt64(photoId))
	if err != nil {
		ctx.Logger.WithError(err).Error("could not get likes list")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return likes list
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(likeList)
}
