package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/egoholic/editor/account/signingin"
)

const (
	getLoginQuery = `SELECT login
	                 FROM accounts
								   WHERE login              = $1
									   AND encrypted_password = $2
								   LIMIT 1;`
	createSigninQuery = `INSERT INTO signins (login, created_at)
																		VALUES ($1,    $2)
											 RETURNING access_token;`
)

type Repository struct {
	db     *sql.DB
	ctx    context.Context
	logger *log.Logger
}

func New(ctx context.Context, db *sql.DB, logger *log.Logger) *Repository {
	return &Repository{db: db, ctx: ctx, logger: logger}
}

var WrongPasswordOrLoginError = errors.New("wrong password or login")

func (r *Repository) IsAuthenticated(login, encryptedPassword string) (bool, error) {
	var (
		err      error
		gotLogin string
	)
	row := r.db.QueryRowContext(r.ctx, login, encryptedPassword)
	err = row.Scan(&gotLogin)
	if err != nil {
		return false, WrongPasswordOrLoginError
	}
	if gotLogin != login {
		return false, WrongPasswordOrLoginError
	}
	return true, nil
}

func (r *Repository) Save(s *signingin.Signin) (at string, err error) {
	row := r.db.QueryRowContext(r.ctx, createSigninQuery, s.Login, s.CreatedAt)
	if err != nil {
		return
	}
	err = row.Scan(&at)
	return
}
