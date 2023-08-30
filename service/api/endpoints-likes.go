package api

import (
    "encoding/json"
    "github.com/bibi2001/WASAPhoto/service/api/reqcontext"
    "github.com/julienschmidt/httprouter"
    "net/http"
)

func (rt *_router) LikesEndpoints(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    // Read the username to be followed from the request body.
    var photoId string
    err := json.NewDecoder(r.Body).Decode(&photoId)
    if err != nil {
        // The body was not a parseable JSON, reject it
        http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
        return
    }

    photoExists, err := rt.db.PhotoExists(photoId)
    if err != nil {
        // In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
        // Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
        ctx.Logger.WithError(err).Error("can't validate photoId")
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    if !photoExists {
        // Here we validate the Request Body photoId, and we
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
	authUser, err := rt.db.GetUsername(token)
	if err != nil {
        ctx.Logger.WithError(err).Error("authenticated username can not be found")
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

	// Check if we should like or unlike the photo
	if r.Method == http.MethodPost {
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
