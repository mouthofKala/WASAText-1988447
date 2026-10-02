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

func (rt *_router) makeChatGroup(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext,
) {
	var request struct {
		ChatName    string   `json:"chatname"`
		Photo       []byte   `json:"photo"`
		Members     []string `json:"members"`
		GrouporChat bool     `json:"grouporchat"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if !(errors.Is(err, nil)) {
		ctx.Logger.WithError(err).Error("error reading request body")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	members := append(request.Members, ctx.UserID)

	if request.Photo != nil && request.GrouporChat == false {
		ctx.Logger.Error("bad req: requesting private chat but attaching a photo")
		http.Error(w, database.BR, http.StatusBadRequest)
	}
	chatID, err := uuid.NewV4()
	if !(errors.Is(err, nil)) {
		ctx.Logger.WithError(err).Error("error generating group or chat ID")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	photouri := database.DefaultGroupPhoto
	basepath := "storage/data"

	err = rt.db.ValidateChatGroup(members, request.GrouporChat)
	if errors.Is(err, database.ErrUserNotFound) {
		ctx.Logger.WithError(err).Error("bad req: one or more users not found")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	if errors.Is(err, database.ErrNoMembersSelected) {
		ctx.Logger.WithError(err).Error("bad req: no members in request")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	if errors.Is(err, database.ErrBadReq) {
		ctx.Logger.WithError(err).Error("bad req: requesting private chat but len(members!=2)")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	if errors.Is(err, database.ErrChatAlreadyExists) {
		ctx.Logger.WithError(err).Error("chat already exists")
		http.Error(w, database.C, http.StatusConflict)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error validating request for group or chat")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	if request.Photo != nil {
		p, err := rt.storage.ParseURI(request.Photo, chatID.String())
		if err != nil {
			ctx.Logger.WithError(err).Error("bad req: invalid image")
			http.Error(w, database.BR, http.StatusBadRequest)
			return
		}
		photouri = basepath + p
	}

	chat, err := rt.db.MakeChatGroup(
		chatID.String(),
		request.ChatName,
		photouri,
		members,
		request.GrouporChat,
	)

	if errors.Is(err, database.ErrNoMembersSelected) {
		ctx.Logger.WithError(err).Error("bad req: no members eligible/found. rolling back...")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	if errors.Is(err, database.ErrBadReq) {
		ctx.Logger.WithError(err).Error("bad req: private chat with more or less than 2 members")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	if errors.Is(err, database.ErrChatAlreadyExists) {
		ctx.Logger.WithError(err).Error("this private chat already exists!")
		http.Error(w, database.BR, http.StatusConflict)
		return
	}
	if errors.Is(err, database.ErrTX) {
		ctx.Logger.WithError(err).Error("error during transaction")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	if errors.Is(err, database.ErrExec) {
		ctx.Logger.WithError(err).Error("error during exec in transaction")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if errors.Is(err, database.ErrGeneratingID) {
		ctx.Logger.WithError(err).Error("error generating new chat or group ID")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("other internal server error")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if request.Photo != nil {
		_, err := rt.storage.SavePFP(request.Photo, chatID.String())
		if errors.Is(err, database.ErrStorage) {
			ctx.Logger.WithError(err).Error("error saving group picture")
			http.Error(w, database.BR, http.StatusInternalServerError)
			return
		}
		if errors.Is(err, database.ErrInvalidImage) {
			ctx.Logger.WithError(err).Error("invalid image")
			http.Error(w, database.BR, http.StatusBadRequest)
			return
		}
		if err != nil {
			ctx.Logger.WithError(err).Error("generic 500 err while saving group picture")
			http.Error(w, database.BR, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(chat); !errors.Is(err, nil) {
		ctx.Logger.WithError(err).Error("error encoding created group")
		return
	}
}
