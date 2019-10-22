package postgresql

import (
	"context"
	"database/sql"
	"log"
)

const (
	accessTokenQuery = `SELECT auth_token
									    FROM accounts
									    WHERE login = $1
								    	AND   encryptedPassword = $2
								    	LIMIT 1;`
)

type Repository struct {
	db     *sql.DB
	ctx    context.Context
	logger *log.Logger
}

func New(ctx context.Context, db *sql.DB, logger *log.Logger) *Repository {
	return &Repository{db: db, ctx: ctx, logger: logger}
}

func (r *Repository) AccessToken(login, encryptedPassword string) (at string, err error) {
	row := r.db.QueryRowContext(r.ctx, accessTokenQuery, login, encryptedPassword)
	err = row.Scan(&at)
	return
}
