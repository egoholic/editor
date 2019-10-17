package postgresql

import (
	"context"
	"database/sql"
	"log"
)

type Repository struct {
	db     *sql.DB
	ctx    context.Context
	logger *log.Logger
}
