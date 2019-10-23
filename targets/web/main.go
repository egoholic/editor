package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	accountSigningin "github.com/egoholic/editor/account/signingin/handler/http"
	. "github.com/egoholic/editor/config"
	rtr "github.com/egoholic/router"
	"github.com/egoholic/router/errback"
	"github.com/egoholic/router/handler"
	"github.com/egoholic/router/node"
	"github.com/egoholic/router/params"
)

type HandlerFnBuilder func(context.Context, *sql.DB, *log.Logger, *errback.ErrBack) func(http.ResponseWriter, *http.Request, *params.Params)

var (
	logger  = log.New(LogFile, "blog", 0)
	connStr string
	db      *sql.DB
	err     error
	errBack *errback.ErrBack
)

func main() {
	errBack, err = errback.New(
		errback.WithBadRequest(badRequestHandler),
		errback.WithNotFound(notFoundHandler),
		errback.WithServerError(serverErrorHandler),
		errback.WithUnauthorized(unauthorizedHandler),
	)
	if err != nil {
		panic(err)
	}
	router := rtr.New()
	root := router.Root()
	signin := root.Child("signin", &node.DumbForm{})
	signin.POST(prepare(accountSigningin.New), "performs sign-in")
}

func prepare(hb HandlerFnBuilder) handler.HandlerFn {
	return func(w http.ResponseWriter, r *http.Request, p *params.Params) {
		d := 100 * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()
		h := hb(ctx, db, logger, errBack)
		h(w, r, p)
	}
}

func badRequestHandler(w http.ResponseWriter, r *http.Request, p *params.Params) {
}
func notFoundHandler(w http.ResponseWriter, r *http.Request, p *params.Params) {
}
func serverErrorHandler(w http.ResponseWriter, r *http.Request, p *params.Params) {
}
func unauthorizedHandler(w http.ResponseWriter, r *http.Request, p *params.Params) {
}
