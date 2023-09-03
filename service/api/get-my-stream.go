package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/bibi2001/WASAPhoto/service/api/reqcontext"
	"github.com/bibi2001/WASAPhoto/service/utils"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the Bearer Token in the header
	token, err := utils.GetBearerToken(r)
	if err != nil {
		http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
		return
	}

	// Get the username corresponding to the Token
	authUser, err := rt.db.GetUsername(token)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("authenticated username cannot be found")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userStream, err := rt.db.GetUserStream(authUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("could not get user stream")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return user stream
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(userStream)

}
