func (rt *_router) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON request body to get the username.
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

	// Return the user identifier in the response
	responseData := LoginResponse{
		Identifier: userIdentifier,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Use http.StatusOK for success
	json.NewEncoder(w).Encode(responseData)
}
