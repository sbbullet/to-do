package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sbbullet/to-do/logger"
	"github.com/sbbullet/to-do/util"
)

type registerUserRequest struct {
	Username string `json:"username" validate:"required,alphanum,min=2,max=12"`
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=16"`
}

func (s *Server) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error(err.Error())
		util.RespondWithBadRequest(w, "Invalid request payload")
	}

	validateRequest(req)
}

func validateRequest(req interface{}) {
	err := validate.Struct(req)
	if err != nil {

		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println(err.Error())
		}
		return
	}
}
