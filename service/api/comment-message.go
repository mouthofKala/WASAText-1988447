package api

//TBD in api handler:
//in filesystem create a dir storage/data/reactions
//and each file has the msgID as the name. json containing a list of pairs: userID, unicode char

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentMessage(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext,
) {
	msgID := ps.ByName("messageID")
	chatID := ps.ByName("chatID")
	userID := ctx.UserID
	reactionID, err := uuid.NewV4()
	if err != nil {
		ctx.Logger.WithError(err).Error("error generating a reactionID")
		return
	}
	var request struct {
		Emoji string `json:"emoji"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ctx.Logger.WithError(err).Error("error decoding request body")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	emojiRunes := []rune(request.Emoji)
	if len(emojiRunes) != 1 || !emojiList[emojiRunes[0]] {
		ctx.Logger.Error("invalid reaction")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}

	reaction, err := rt.db.CommentMessage(msgID, chatID, userID, reactionID.String(), emojiRunes[0])
	if errors.Is(err, database.ErrBadReq) {
		ctx.Logger.WithError(err).Error("bad request: chatID or messageID do not exist")
		http.Error(w, database.BR, http.StatusBadRequest)
		return
	}
	if errors.Is(err, database.ErrForbidden) {
		ctx.Logger.WithError(err).Error("forbidden: userid not member of chatid")
		http.Error(w, database.F, http.StatusForbidden)
		return
	}
	if errors.Is(err, database.ErrExec) {
		ctx.Logger.WithError(err).Error("500 coming from a db.exec")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("generic 500")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	if err = rt.storage.AddReaction(reaction, string(chatID)); err != nil {
		ctx.Logger.WithError(err).Error("error saving reaction")
		http.Error(w, database.ISE, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reaction); err != nil {
		ctx.Logger.WithError(err).Error("error encoding reaction response")
		return
	}

}
