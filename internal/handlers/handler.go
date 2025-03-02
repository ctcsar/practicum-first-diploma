package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ctcsar/practicum-first-diploma/internal/auth"
	"github.com/ctcsar/practicum-first-diploma/internal/compress"
	"github.com/ctcsar/practicum-first-diploma/internal/logger"
	"github.com/ctcsar/practicum-first-diploma/internal/storage"
)

type userData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var buff userData
	var authUser *auth.User
	var token string

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&buff)
	if err != nil {
		logger.Log.Error("cannot decode user data", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if buff.Login == "" || buff.Password == "" {
		logger.Log.Error("login or password is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authUser, err = auth.CreateUser(buff.Login, buff.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Log.Error("cannot create user", zap.Error(err))
		return
	}
	err = storage.SaveUser(db, authUser)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, "User with this login already exists")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log.Error("cannot save user", zap.Error(err))
		}
		return
	}

	token, err = auth.AuthUser(buff.Login, buff.Password, db)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		logger.Log.Error("cannot auth user", zap.Error(err))
		return
	}

	w.Header().Set("Authorization", token)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, token)
}

func LoginUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var buff userData
	token := r.Header.Get("Authorization")

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&buff)
	if err != nil {
		logger.Log.Error("cannot decode user data", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if buff.Login == "" || buff.Password == "" {
		logger.Log.Error("login or password is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if token != "" {
		_, err := auth.VerifyToken(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log.Error("cannot verify token", zap.Error(err))
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "User already logged in")
		return

	}

	user, err := auth.AuthUser(buff.Login, buff.Password, db)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		logger.Log.Error("cannot auth user", zap.Error(err))
		return
	}
	w.Header().Set("Authorization", user)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, user)
}

func Routers(ctx context.Context, handler chi.Router, db *sql.DB) {
	handler.Post("/api/user/register", func(w http.ResponseWriter, r *http.Request) {
		RegisterUser(w, r, db)
	})
	handler.Post("/api/user/login", func(w http.ResponseWriter, r *http.Request) {
		LoginUser(w, r, db)
	})
}

func Run(ctx context.Context, url string, handler chi.Router, db *sql.DB) error {
	logger.Log.Info("starting server", zap.String("url", url))
	handler = logger.RequestLogger(handler)
	Routers(ctx, handler, db)
	return http.ListenAndServe(url, compress.GzipMiddleware(handler))
}
