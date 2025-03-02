package storage

import (
	"database/sql"

	"github.com/ctcsar/practicum-first-diploma/internal/auth"
	"github.com/ctcsar/practicum-first-diploma/internal/database"
)

func SaveUser(db *sql.DB, user *auth.User) error {
	err := database.SaveUser(db, user.ID, user.Login, user.Password)
	return err
}
