package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) addToGroup(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext,
) {
	var request struct {
		TUserID string `json:"targetuserID"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ctx.Logger.WithError(err).Error("error reading request")
		http.Error(w, database.BR, http.StatusBadRequest)
	}
	userID := ctx.UserID
	GroupID := ps.ByName("chatID")
	chat, err := rt.db.AddToGroup(request.TUserID, GroupID, userID)
	//list off all errors with logger and http
	if err == database.ErrUserNotFound {
		ctx.Logger.WithError(err).Error("forbidden: requestor not in group!")
		http.Error(w, database.F, http.StatusForbidden)
		return
	}
	if err == database.ErrNotaGroup {
		ctx.Logger.WithError(err).Error("forbidden: can't add members to a private chat")
		http.Error(w, database.F, http.StatusForbidden)
		return
	}
	if err == database.ErrBadReq {
		ctx.Logger.WithError(err).Error("bad request: no group found")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}

	if err == database.ErrUserAlreadyIn {
		ctx.Logger.WithError(err).Error("bad request: targetuser already in group")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}

	if err == database.ErrModifyingChat {
		ctx.Logger.WithError(err).Error("Internal server error: error modifying chat")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if err == database.ErrFetchingChat {
		ctx.Logger.WithError(err).Error("internal server error: error compiling new chat json")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if err == database.ErrFetchingMembers {
		ctx.Logger.WithError(err).Error("internal server err: error fetching members to compile chat json")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if err != nil {
		ctx.Logger.WithError(err).Error("other internal server error")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(chat); err != nil {
		ctx.Logger.WithError(err).Error("error encoding updated chat")
		return
	}

}
