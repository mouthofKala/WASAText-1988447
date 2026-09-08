package api

import (
	"database/sql"
	"errors"
	"io"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setGroupPhoto(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext) {

	ChatID := ps.ByName("chatID")
	userID := ctx.UserID

	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	newphotobytes, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("error reading uploaded photo")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}

	//save newphotobytes to a uri
	photoURI, err := rt.storage.SavePFP(newphotobytes, ChatID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error storing photo")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	oldphotoURI, err := rt.db.SetGroupPhoto(photoURI, ChatID, userID)

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
		ctx.Logger.WithError(err).Error("error setting photo. rolling back...")
		_ = rt.storage.DeletePFP(photoURI)
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if err := rt.storage.DeletePFP(oldphotoURI); err != nil {
		ctx.Logger.WithError(err).Error("error deleting old group photo")
	}
	w.WriteHeader(http.StatusNoContent)
}
