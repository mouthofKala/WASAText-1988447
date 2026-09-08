package api

import (
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

	//FOR THIS TO WORK THE FRONTEND HAS TO HAVE AN IMG REQUESTBODY
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	newphotobytes, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("error reading uploaded photo")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}

	//save newphotobytes to a uri
	photoURI, err := rt.storage.SavePFP(newphotobytes, string(ctx.UserID))
	if err != nil {
		ctx.Logger.WithError(err).Error("error storing photo")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	oldphotouri, err := rt.db.SetMyPhoto(photoURI, string(ctx.UserID))
	if err != nil {
		ctx.Logger.WithError(err).Error("error setting photo. rolling back...")
		_ = rt.storage.DeletePFP(photoURI)
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	if err = rt.storage.DeletePFP(oldphotouri); err != nil {
		ctx.Logger.WithError(err).Error("error deleting old pfp")
	}

	w.WriteHeader(http.StatusNoContent)
}
