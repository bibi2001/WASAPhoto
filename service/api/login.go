import (
	"github.com/julienschmidt/httprouter"
	"github.com/bibi2001/WASAPhoto/service/api/reqcontext"
	"net/http"
	"encoding/json"
	"net/http"
)

func (rt *_router) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parsing the JSON request body to get the username
	var requestBody struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username := requestBody.Name
	userIdentifier, err := rt.db.GetUserId(username)
	if err != nil {
		http.Error(w, "Failed to retrieve user identifier", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// Set the Bearer token
    w.Header().Set("Authorization", "Bearer "+userIdentifier) 
    w.WriteHeader(http.StatusCreated)

	// Return the user identifier in the response
	responseData := LoginResponse{
		Identifier: userIdentifier,
	}
	json.NewEncoder(w).Encode(responseData)
}
