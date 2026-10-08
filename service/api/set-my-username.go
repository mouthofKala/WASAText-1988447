package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// read requestbody json to ckeck 400 bad request
	var request struct {
		Username string `json:"targetusername"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error decoding request body")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}

	validUsername := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !validUsername.MatchString(request.Username) || len(request.Username) < 3 || len(request.Username) > 16 {
		ctx.Logger.WithError(err).Error("username length and letter constrains violated")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}

	err = rt.db.SetMyUsername(request.Username, ctx.UserID)
	if !errors.Is(err, nil) {
		if errors.Is(err, database.ErrUsernameUnavailable) {
			ctx.Logger.WithError(err).Error("conflict: username already in use")
			http.Error(w, database.C, http.StatusConflict)
			return
		}

		ctx.Logger.WithError(err).Error("error setting username")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
