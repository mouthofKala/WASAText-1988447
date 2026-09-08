package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) makeGroup(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext,
) {
	var request struct {
		ChatName string   `json:"chatName"`
		Photo    string   `json:"photo"`
		Members  []string `json:"members"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ctx.Logger.WithError(err).Error("error reading request body")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	members := append(request.Members, ctx.UserID)
	chat, err := rt.db.MakeGroup(
		request.ChatName,
		request.Photo,
		members,
	)

	if errors.Is(err, database.ErrNoMembersSelected) {
		ctx.Logger.WithError(err).Error("no members eligible/found")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}

	if err != nil {
		ctx.Logger.WithError(err).Error("error creating chat")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(chat); err != nil {
		ctx.Logger.WithError(err).Error("error encoding created group")
		return
	}
}
