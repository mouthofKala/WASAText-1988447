package api

import (
	"database/sql"
	"encoding/json"
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

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
	newphotobytes, err := io.ReadAll(r.Body)
	if !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error reading uploaded photo")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}

	// save newphotobytes to a uri
	photoURI, err := rt.storage.SavePFP(newphotobytes, ChatID)
	if !errors.Is(err, nil) {
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

	if !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error setting photo. rolling back...")
		_ = rt.storage.DeletePFP(photoURI)
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if err := rt.storage.DeletePFP(oldphotoURI); !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error deleting old group photo")
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(photoURI); !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error encoding new group photouri")
	}
}
