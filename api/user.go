package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/sbbullet/to-do/db"
	"github.com/sbbullet/to-do/logger"
	"github.com/sbbullet/to-do/util"
)

type userResponse struct {
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type registerUserRequest struct {
	Username string `json:"username" validate:"required,alphanum,min=2,max=12"`
	FullName string `json:"full_name" validate:"required,full_name"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=16"`
}

func createUserResponse(user db.User) userResponse {
	return userResponse{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (s *Server) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.RespondWithBadRequest(w, "Invalid request payload")
		return
	}

	validationErrors := validateRequest(req)
	if validationErrors != nil {
		util.RespondWithValidationErrors(w, validationErrors)
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		logger.Error(err.Error())
		util.RespondWithInternalServerError(w)
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		Email:          req.Email,
		FullName:       req.FullName,
		HashedPassword: hashedPassword,
	}

	user, err := s.store.CreateUser(arg)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			field := strings.Split(strings.SplitN(err.Error(), ":", 2)[1], ".")[1]
			validationErrors := map[string][]string{
				field: {"This " + field + " is already taken"},
			}
			util.RespondWithValidationErrors(w, validationErrors)
			return
		}
		logger.Error(err.Error())
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithOk(w, createUserResponse(user))
}

type loginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginUserResponse struct {
	AccessToken          string       `json:"access_token"`
	AccessTokenExpiresAt time.Time    `json:"access_token_expires_at"`
	User                 userResponse `json:"user"`
}

func (s *Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req loginUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.RespondWithBadRequest(w, "Invalid request payload")
		return
	}

	validationErrors := validateRequest(req)
	if validationErrors != nil {
		util.RespondWithValidationErrors(w, validationErrors)
		return
	}

	user, err := s.store.GetUser(req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.RespondWithUauthorizedError(w, "Oops!! These credentials do not match our records")
			return
		}
		logger.Error(err.Error())
		util.RespondWithInternalServerError(w)
		return
	}

	if err := util.CheckHashPassword(user.HashedPassword, req.Password); err != nil {
		util.RespondWithUauthorizedError(w, "Oops!! These credentials do not match our records")
		return
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		logger.Error(err.Error())
		util.RespondWithInternalServerError(w)
		return
	}

	response := loginUserResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiresAt,
		User:                 createUserResponse(user),
	}

	util.RespondWithOk(w, response)
}
