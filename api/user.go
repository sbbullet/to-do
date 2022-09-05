package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sbbullet/to-do/db"
	"github.com/sbbullet/to-do/logger"
	"github.com/sbbullet/to-do/util"
)

type userResponse struct {
	Username string `json:"username" validate:"required,alphanum,min=2,max=12"`
	FullName string `json:"full_name" validate:"required,full_name"`
	Email    string `json:"email" validate:"required,email"`
}

type registerUserRequest struct {
	Username string `json:"username" validate:"required,alphanum,min=2,max=12"`
	FullName string `json:"full_name" validate:"required,full_name"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=16"`
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

	util.RespondWithOk(w, user)
}

type loginUserRequest struct {
}
