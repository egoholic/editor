package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/egoholic/editor/account/signingin"
	repository "github.com/egoholic/editor/account/signingin/repository/postgresql"
	"github.com/egoholic/router/errback"
	"github.com/egoholic/router/params"
)

type form struct {
	Login    []byte `json:"login"`
	Password []byte `json:"password"`
}
type responseData struct {
	AccessToken string `json:"accessToken"`
}

func New(ctx context.Context, db *sql.DB, logger *log.Logger, errBack *errback.ErrBack) func(w http.ResponseWriter, r *http.Request, p *params.Params) {
	return func(w http.ResponseWriter, r *http.Request, p *params.Params) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errBack.HandleServerError(w, r, p)
			return
		}
		form := &form{}
		err = json.Unmarshal(body, form)
		if err != nil {
			errBack.HandleServerError(w, r, p)
			return
		}
		repo := repository.New(ctx, db, logger)
		value, err := signingin.New(logger, repo, form.Login, form.Password)
		if err != nil {
			errBack.HandleNotFound(w, r, p)
			return
		}
		w.WriteHeader(200)
		rd := &responseData{
			AccessToken: value.AccessToken(),
		}
		bytes, err := json.Marshal(rd)
		if err != nil {
			errBack.HandleServerError(w, r, p)
			return
		}
		w.Write(bytes)
	}
}
