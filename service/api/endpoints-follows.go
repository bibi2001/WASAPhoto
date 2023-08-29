package api

import (
	"encoding/json"
	"github.com/bibi2001/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Read the username to be followed from the request body.
	var username string
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		// The body was not a parseable JSON, reject it
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create the fountain in the database. Note that this function will return a new instance of the fountain with the
	// same information, plus the ID.
	dbfountain, err := rt.db.FollowUser(, username)
	if err != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		ctx.Logger.WithError(err).Error("can't create the fountain")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Here we can re-use `fountain` as FromDatabase is overwriting every variabile in the structure.
	fountain.FromDatabase(dbfountain)

	// Send the output to the user.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(fountain)
}
