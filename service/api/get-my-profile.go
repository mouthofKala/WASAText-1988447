package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyProfile(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext,
) {
	user, err := rt.db.GetUserProfile(ctx.UserID)

	if errors.Is(err, database.ErrUserNotFound) {
		ctx.Logger.WithError(err).Error("user not found")
		http.Error(w, database.NF, http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		ctx.Logger.WithError(err).Error("error encoding profile")
		return
	}

}
