package api

import (
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) leaveGroup(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext,
) {
	userID := ctx.UserID
	chatID := ps.ByName("chatID")

	lastmember, err := rt.db.LeaveGroup(userID, chatID)

	if errors.Is(err, database.ErrBadReq) {
		ctx.Logger.WithError(err).Error("bad request: chat/group does not exist")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	if errors.Is(err, database.ErrNotaGroupMember) {
		ctx.Logger.WithError(err).Error("forbidden: user is not a member!")
		http.Error(w, database.F, http.StatusForbidden)
		return
	}
	if errors.Is(err, database.ErrFetchingMembers) {
		ctx.Logger.WithError(err).Error("Error fetching members of chat/group")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	if errors.Is(err, database.ErrNoMembersSelected) {
		ctx.Logger.WithError(err).Error("orphaned chat/group or I.S.E")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	if errors.Is(err, database.ErrTX) {
		ctx.Logger.WithError(err).Error("ERROR DURING A TRANSACTION")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if err != nil {
		ctx.Logger.WithError(err).Error("generic 500")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if lastmember {
		err = rt.storage.DeleteChatGroup(chatID)
		if errors.Is(err, database.ErrStorage) {
			ctx.Logger.WithError(err).Error("error deleting group or chat items from filesys")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)

}
