package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyPhoto(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext) {

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
	newphotobytes, err := io.ReadAll(r.Body)
	if !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error reading uploaded photo")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}

	// save newphotobytes to a uri
	photoURI, err := rt.storage.SavePFP(newphotobytes, ctx.UserID)
	if !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error storing photo")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	oldphotouri, err := rt.db.SetMyPhoto(photoURI, ctx.UserID)
	if !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error setting photo. rolling back...")
		if err := rt.storage.DeletePFP(photoURI); !errors.Is(err, nil) {
			ctx.Logger.WithError(err).Error("error rolling back pfp")
		}
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	if err := rt.storage.DeletePFP(oldphotouri); !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error deleting old pfp")
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(photoURI); !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("Error encoding new photouri")
	}
}
