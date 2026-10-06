package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID := ctx.UserID
	chats, err := rt.db.GetMyConversations(userID)
	if !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error fetching chats")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	/*
		if len(chats) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	*/

	// now read chats and send them in json
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(chats); !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error encoding chats")
		return
	}
}
