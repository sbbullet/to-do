package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sbbullet/to-do/db"
	"github.com/sbbullet/to-do/logger"
	"github.com/sbbullet/to-do/util"
)

type createTodoRequest struct {
	Title string `json:"title" validate:"required,min=6,max=255"`
}

// Create todo for the authorized user
func (s *Server) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithBadRequest(w, "Invalid request payload")
		return
	}

	validationErrors := validateRequest(req)
	if validationErrors != nil {
		util.RespondWithValidationErrors(w, validationErrors)
		return
	}

	todoID, err := uuid.NewRandom()
	if err != nil {
		logger.Error(err.Error())
		util.RespondWithInternalServerError(w)
		return
	}

	username := r.Header.Get(authUsernameHeaderKey)

	arg := db.CreateTodoParams{
		ID:       todoID,
		Username: username,
		Title:    req.Title,
	}

	todo, err := s.store.CreateTodo(arg)
	if err != nil {
		logger.Error(err.Error())
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithOk(w, todo)
}

// Get todos of the authorized user
func (s *Server) GetUserTodos(w http.ResponseWriter, r *http.Request) {
	var pageNum int
	var pageSize int

	pageNum, err := strconv.Atoi(r.URL.Query().Get("page_num"))
	if err != nil {
		pageNum = 1
	}

	pageSize, err = strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 5
	}

	if pageNum <= 0 || pageSize <= 0 {
		util.RespondWithBadRequest(w, "Page number and page size must be greater than zero")
		return
	}

	username := r.Header.Get(authUsernameHeaderKey)

	arg := db.GetUserTodosParams{
		Username: username,
		Limit:    pageSize,
		Offset:   (pageNum - 1) * pageSize,
	}

	todos, err := s.store.GetUserTodos(arg)
	fmt.Println(todos)
	if err != nil {
		logger.Error(err.Error())
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithOk(w, createTodosResponse(todos))
}

type updateTodoRequest struct {
	Title       string `json:"title" validation:"min=6,max=255"`
	IsCompleted *bool  `json:"is_completed" validation:"boolean"`
}

// Update specified todo of the authorized user
func (s *Server) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var req updateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err.Error())
		util.RespondWithBadRequest(w, "Invalid request payload")
		return
	}

	todoID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		util.RespondWithBadRequest(w, "Invalid todo identifier")
		return
	}

	validationErrors := validateRequest(req)
	if validationErrors != nil {
		util.RespondWithValidationErrors(w, validationErrors)
		return
	}

	username := r.Header.Get(authUsernameHeaderKey)
	todo, err := s.store.GetTodoById(todoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.RespondWithNotFoundError(w, "Oops!! We couldn't find the associated todo")
			return
		}
	}

	if todo.Username != username {
		util.RespondWithForbiddenError(w, "You are forbidden to perform the action on this resource")
		return
	}

	var isCompleted sql.NullBool
	if req.IsCompleted != nil {
		isCompleted = sql.NullBool{Bool: *req.IsCompleted, Valid: true}
	} else {
		isCompleted = sql.NullBool{Valid: false}
	}

	updateTodoArgs := db.UpdateTodoParams{
		ID:          todo.ID,
		Title:       sql.NullString{String: req.Title, Valid: len(req.Title) > 0},
		IsCompleted: isCompleted,
	}
	updatedTodo, err := s.store.UpdateTodo(updateTodoArgs)
	if err != nil {
		logger.Error(err.Error())
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithOk(w, createTodoResponse(updatedTodo))
}

// Delete specified todo of the authorized user
func (s *Server) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		util.RespondWithBadRequest(w, "Invalid todo identifier")
		return
	}

	username := r.Header.Get(authUsernameHeaderKey)
	arg := db.DeleteTodoOfAUserParams{
		ID:       todoId,
		Username: username,
	}

	err = s.store.DeleteTodoOfAUser(arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.RespondWithNotFoundError(w, "Oops! We couldn't any of your todos with given identifier")
			return
		}

		logger.Error(err.Error())
		util.RespondWithInternalServerError(w)
		return
	}

	util.RespondWithOk(w, "Successfully deleted specified todo from your todo list")
}

type todoResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	CreatedAt   time.Time `json:"created_at"`
	IsCompleted bool      `json:"is_completed"`
}

func createTodoResponse(todo db.Todo) todoResponse {
	return todoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		IsCompleted: todo.IsCompleted,
		CreatedAt:   todo.CreatedAt,
	}
}

func createTodosResponse(todos []db.Todo) []todoResponse {
	todosToSend := []todoResponse{}

	for _, todo := range todos {
		todosToSend = append(todosToSend, createTodoResponse(todo))
	}

	return todosToSend
}
