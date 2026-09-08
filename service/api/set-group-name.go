package api

import (
	"database/sql"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setGroupName(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext) {

	ChatID := ps.ByName("chatID")
	userID := ctx.UserID
	var request struct {
		NewName string `json:"newname"`
	}

	err := rt.db.SetGroupName(request.NewName, ChatID, userID)

	if errors.Is(err, database.ErrNotaGroup) {
		ctx.Logger.WithError(err).Error("can't change private chat pfps")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}

	if errors.Is(err, database.ErrNotaGroupMember) {
		ctx.Logger.WithError(err).Error("Forbidden: not a group member")
		http.Error(w, database.F, http.StatusForbidden)
		return
	}

	if errors.Is(err, sql.ErrNoRows) {
		ctx.Logger.WithError(err).Error("group not found")
		http.Error(w, database.NF, http.StatusNotFound)
		return
	}

	if err != nil {
		ctx.Logger.WithError(err).Error("error setting photo")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
