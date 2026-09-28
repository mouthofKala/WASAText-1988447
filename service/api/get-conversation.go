package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID := ctx.UserID
	chatID := ps.ByName("chatID")

	convo, err := rt.db.GetConversation(userID, chatID)
	if err == database.ErrFetchingChat {
		ctx.Logger.WithError(err).Error("chat not found or user not authorised")
		http.Error(w, database.NF, http.StatusNotFound)
		return
	}

	if err == database.ErrFetchingMembers {
		ctx.Logger.WithError(err).Error("error with fetching or scanning members from db")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if err == database.ErrFetchingMessages {
		ctx.Logger.WithError(err).Error("error with fetching or scanning messages from db")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if err != nil {
		ctx.Logger.WithError(err).Error("generic 500")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(convo); err != nil {
		ctx.Logger.WithError(err).Error("error encoding conversation response")
	}
	//but it succeeded so no need to send http error?

}
