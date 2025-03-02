package database

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/google/uuid"

	"time"
)

type DBUser struct {
	ID           uuid.UUID `json:"id"`
	Login        string    `json:"login"`
	Password     string    `json:"password"`
	RegisterDate string    `json:"registerDate,omitempty"`
}

func SaveUser(db *sql.DB, id uuid.UUID, login string, password string) error {
	status, err := db.Query("INSERT INTO users (id, login, password, registerDate) VALUES ($1, $2, $3, $4)", id, login, password, time.Now())
	if err != nil {
		pgError, ok := err.(*pq.Error)
		if ok {
			// Check the PostgreSQL error code
			switch pgError.Code {
			case "23505": // unique_violation
				return fmt.Errorf("user with login %s already exists", login)
			case "23503": // foreign_key_violation
				return fmt.Errorf("foreign key constraint failed")
			default:
				return fmt.Errorf("database error: %v", pgError)
			}
		}
		return err
	}
	defer status.Close()
	return nil
}

func CheckUser(login string, password string, db *sql.DB) (uuid.UUID, string, error) {
	var user DBUser
	fmt.Printf("Checking user with password %s\n", password)
	result := db.QueryRow("SELECT * FROM users WHERE login = $1", login)
	err := result.Scan(&user.ID, &user.Login, &user.Password, &user.RegisterDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, "", fmt.Errorf("user with login %s not found", login)
		}
		return uuid.Nil, "", fmt.Errorf("database error: %v", err)
	}
	fmt.Printf("User with login %s found\n", user)
	return user.ID, user.Password, nil
}
