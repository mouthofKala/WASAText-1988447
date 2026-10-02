package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext,
) {
	userID := ctx.UserID
	chatID := ps.ByName("chatID")

	msgid, err := uuid.NewV4()
	if !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error generating messgae ID")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(5 << 20); !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error parsing multipart req")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}

	var content *string
	if value := r.FormValue("content"); value != "" {
		content = &value
	}

	var replyto *string
	if value := r.FormValue("replyto"); value != "" {
		replyto = &value
	}

	var fwdfrom *string
	if value := r.FormValue("fwdfrom"); value != "" {
		fwdfrom = &value
	}

	msgID := msgid.String()

	var photouri *string
	file, _, err := r.FormFile("photo")
	if errors.Is(err, nil) {
		defer file.Close()
		data, err := io.ReadAll(file)
		if !errors.Is(err, nil) {
			ctx.Logger.WithError(err).Error("error reading msg attachment")
			http.Error(w, database.BR, http.StatusBadRequest)
			return
		}
		uri, err := rt.storage.SaveMSGPhoto(data, msgID, chatID)
		if !errors.Is(err, nil) {
			ctx.Logger.WithError(err).Error("error storing msg attachment")
			if errors.Is(err, database.ErrInvalidImage) {
				http.Error(w, database.BR, http.StatusBadRequest)
			} else {
				http.Error(w, database.ISE, http.StatusInternalServerError)
			}
			return
		}
		photouri = &uri
	} else if !errors.Is(err, http.ErrMissingFile) {
		ctx.Logger.WithError(err).Error("error fetching msg attachment")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}

	msg, err := rt.db.SendMessage(userID, chatID, content, photouri, replyto, fwdfrom, msgID)
	if errors.Is(err, database.ErrBadReq) {

		_ = rt.storage.DeleteMSGFiles(msgID, chatID)
		ctx.Logger.WithError(err).Error("bad request: sending to inexistent chat?")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	if errors.Is(err, database.ErrNotaGroupMember) {
		_ = rt.storage.DeleteMSGFiles(msgID, chatID)
		ctx.Logger.WithError(err).Error("403: not a group member")
		http.Error(w, database.F, http.StatusForbidden)
		return
	}
	if errors.Is(err, database.ErrInvalidMsg) {
		_ = rt.storage.DeleteMSGFiles(msgID, chatID)
		ctx.Logger.WithError(err).Error("400: message empty or invalid")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	if errors.Is(err, database.ErrGeneratingID) {
		_ = rt.storage.DeleteMSGFiles(msgID, chatID)
		ctx.Logger.WithError(err).Error("500: error gneerating ID")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	if errors.Is(err, database.ErrExec) || errors.Is(err, database.ErrTX) {
		_ = rt.storage.DeleteMSGFiles(msgID, chatID)
		ctx.Logger.WithError(err).Error("500: error doing a database Exec or during the transaction")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	if !errors.Is(err, nil) {
		_ = rt.storage.DeleteMSGFiles(msgID, chatID)
		ctx.Logger.WithError(err).Error("generic 500")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	// save json
	if err := rt.storage.SaveMSGJson(msg); !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error saving message JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(msg); !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("Error encoding response")
	}

}
