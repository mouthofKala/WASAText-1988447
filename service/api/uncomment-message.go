package api

import (
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentMessage(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext) {

	msgID := ps.ByName("messageID")
	chatID := ps.ByName("chatID")
	reactionID := ps.ByName("reactionID")
	userID := ctx.UserID

	err := rt.db.UncommentMessage(msgID, chatID, userID, reactionID)

	if errors.Is(err, database.ErrForbidden) {
		ctx.Logger.WithError(err).Error("forbidden: user has no auth")
		http.Error(w, database.F, http.StatusForbidden)
		return
	}

	if errors.Is(err, database.ErrBadReq) {
		ctx.Logger.WithError(err).Error("Bad request: message, chat or reaction nonexisting")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}

	if errors.Is(err, database.ErrExec) {
		ctx.Logger.WithError(err).Error("error in exec on db")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("generic 500")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	err = rt.storage.DeleteReaction(msgID, userID, chatID, reactionID)

	if !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("ERROR DELETING REACTION FROM FILESYS!")
		// you need a function which checks every x time for any orphaned chats, messages or reactions?

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

}
