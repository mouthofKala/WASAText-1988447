package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) makeChat(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext,
) {
	var request struct {
		ChatName string `json:"otheruser"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ctx.Logger.WithError(err).Error("error reading request body")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	creatorID := ctx.UserID

	chat, err := rt.db.MakeChat(
		request.ChatName,
		creatorID,
	)
	if errors.Is(err, database.ErrChatAlreadyExists) {
		ctx.Logger.WithError(err).Error("conflict: private chat already exists")
		http.Error(w, database.C, http.StatusConflict)
		return
	}
	if errors.Is(err, database.ErrGeneratingID) || err != nil {
		ctx.Logger.WithError(err).Error("error generating chat ID")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(chat); err != nil {
		ctx.Logger.WithError(err).Error("error encoding created chat")
		return
	}

}
