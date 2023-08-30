package api

import (
    "encoding/json"
    "github.com/bibi2001/WASAPhoto/service/api/reqcontext"
    "github.com/julienschmidt/httprouter"
    "net/http"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    // Read the username to be followed from the request body.
    var usernameToFollow string
    err := json.NewDecoder(r.Body).Decode(&usernameToFollow)
    if err != nil {
        // The body was not a parseable JSON, reject it
        http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
        return
    }

    usernameExists, err := rt.db.UserExists(usernameToFollow)
    if err != nil {
        // In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
        // Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
        ctx.Logger.WithError(err).Error("can't validate username")
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    if !usernameExists {
        // Here we validate the Request Body username, and we
        // discovered that it is not valid.
        http.Error(w, "Invalid Username", http.StatusBadRequest)
        return
    }

    // Get the Bearer Token in the header
    authUserID, err := GetBearerToken(r)
    if err != nil {
        http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
        return
    }

    err = rt.db.FollowUser(authUserID, usernameToFollow)
    if err != nil {
        ctx.Logger.WithError(err).Error("can't follow user")
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}
