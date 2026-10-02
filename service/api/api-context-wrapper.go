package api

import (
	"net/http"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// httpRouterHandler is used for functions that accept a reqcontext.
// RequestContext is used to add to the context params
// required by the httprouter package.

type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
// the auth generating process mirrors doLogin and is only performed if the function needs authentication

func (rt *_router) wrapPublic(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return rt.wrapAuth(fn, true)
}

func (rt *_router) wrap(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return rt.wrapAuth(fn, false)
}

func (rt *_router) wrapAuth(fn httpRouterHandler, public bool) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			http.Error(w, database.ISE, http.StatusInternalServerError)
			// w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		if !public {
			auth := r.Header.Get("Authorization")
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
				ctx.Logger.Error("malformed or missing Authorization header")
				http.Error(w, database.Un, http.StatusUnauthorized)
				return
			}
			ctx.UserID = parts[1]
			exists, err := rt.db.UserExists(ctx.UserID)
			if err != nil {
				ctx.Logger.WithError(err).Error("error in checking authentication")
				http.Error(w, database.ISE, http.StatusInternalServerError)
				return
			}
			if !exists {
				ctx.Logger.WithError(err).Error("this user does not exist in the db")
				http.Error(w, database.Un, http.StatusUnauthorized)
			}

			// if public, do nothing else
		}

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}
