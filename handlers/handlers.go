package handlers

import "database/sql"

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

type Handler struct {
	db *sql.DB
}
