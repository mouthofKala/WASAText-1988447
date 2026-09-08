package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext,
) {
	var request struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)

	if len(request.Name) < 3 || len(request.Name) > 16 {
		ctx.Logger.WithError(err).Error("invalid username")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("generic bad request")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}

	userID, err := rt.db.DoLogin(request.Name)
	if err != nil {
		ctx.Logger.WithError(err).Error("error during login")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	response := struct {
		Identifier string `json:"identifier"`
	}{
		Identifier: userID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("error encoding login response")
	}
}
