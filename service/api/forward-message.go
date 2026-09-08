package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) forwardMessage(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext,
) {
	userID := ctx.UserID
	chatID := ps.ByName("chatID")
	fwdmessageID := ps.ByName("messageID")
	newmsgID, err := uuid.NewV4()
	if err != nil {
		ctx.Logger.WithError(err).Error("error generating messgae ID")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	var request struct {
		TargetChatID string `json:"targetchat_id"`
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ctx.Logger.WithError(err).Error(database.ErrDecodingReq)
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	msg, err := rt.db.ForwardMessage(userID, request.TargetChatID, fwdmessageID, chatID, newmsgID.String())
	if errors.Is(err, database.ErrForbidden) {
		ctx.Logger.WithError(err).Error("403: forbidden in fwdmsg method")
		http.Error(w, database.F, http.StatusForbidden)
		return
	}
	if errors.Is(err, database.ErrNotaGroupMember) {
		ctx.Logger.WithError(err).Error("403: not a group member")
		http.Error(w, database.F, http.StatusForbidden)
		return
	}
	if errors.Is(err, database.ErrInvalidMsg) {
		ctx.Logger.WithError(err).Error("400: message empty, missing or invalid")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	if errors.Is(err, database.ErrBadReq) {
		ctx.Logger.WithError(err).Error("400: requested msgID not found")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	if errors.Is(err, database.ErrGeneratingID) {
		ctx.Logger.WithError(err).Error("500: error gneerating ID")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	if errors.Is(err, database.ErrExec) {
		ctx.Logger.WithError(err).Error("500: error doing a database Exec")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("generic 500")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if err := rt.storage.SaveMSGJson(msg); err != nil {
		ctx.Logger.WithError(err).Error("error saving message JSON")
	}
	//ADD SAVEMSGPHOTO!!!

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(msg); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
	}

}
