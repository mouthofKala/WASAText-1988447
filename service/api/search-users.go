package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	searchkey := r.URL.Query().Get("searchkey")

	users, err := rt.db.SearchUsers(searchkey)
	if !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error fetching userlist")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	if len(users) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(users); !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error encoding response")
	}
}
