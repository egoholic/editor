package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	logger.Println("server starting...")
	db, err = sql.Open("postgres", DBConnectionString)
	if err != nil {
		logger.Printf("ERROR: %s\n", err.Error())
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logger.Printf("ERROR: %s\n", err.Error())
	}
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
	pid := os.Getpid()
	pidf, err := os.Create(PIDFilePath)
	if err != nil {
		logger.Printf("FATALPID: %s   = %d\n", err.Error(), pid)
	}
	_, err = pidf.WriteString(strconv.Itoa(pid))
	if err != nil {
		logger.Printf("FATALPID: %s   = %d\n", err.Error(), pid)
	}
	defer func() {
		err = pidf.Close()
		if err != nil {
			logger.Printf("FATALPID: %s   = %d\n", err.Error(), pid)
		}
		err = os.Remove(PIDFilePath)
		if err != nil {
			logger.Printf("FATALPID: %s   = %d\n", err.Error(), pid)
		}
	}()

	logger.Printf("server listens :%d port\n", Port)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), router))
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
