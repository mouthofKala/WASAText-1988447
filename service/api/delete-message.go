package api

import (
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteMessage(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext) {
	msgID := ps.ByName("messageID")
	chatID := ps.ByName("chatID")
	userID := ctx.UserID

	err := rt.db.DeleteMessage(msgID, userID, chatID)
	if errors.Is(err, database.ErrBadReq) {
		ctx.Logger.WithError(err).Error("bad request: message does not exist")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	if errors.Is(err, database.ErrForbidden) {
		ctx.Logger.WithError(err).Error("forbidden: user has no auth to delete this msgID")
		http.Error(w, database.F, http.StatusForbidden)
		return
	}
	if errors.Is(err, database.ErrTX) {
		ctx.Logger.WithError(err).Error("error during a go transaction")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	if errors.Is(err, database.ErrExec) {
		ctx.Logger.WithError(err).Error("error during an exec in the transaction")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	if !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("generic 500")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	err = rt.storage.DeleteMSGFiles(msgID, chatID)
	if !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error deleting msgfiles from filesystem")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
