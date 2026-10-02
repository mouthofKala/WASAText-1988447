package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID := ps.ByName("userID")
	user, err := rt.db.GetUserProfile(userID)

	if errors.Is(err, sql.ErrNoRows) {
		ctx.Logger.WithError(err).Error("user not found")
		http.Error(w, database.NF, http.StatusNotFound)
		return
	}

	if !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error fetching user")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error encoding user")
		// http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
