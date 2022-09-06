package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
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
		ID:       todoID.String(),
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

// Update todo of the authorized user

type todoResponse struct {
	ID          string    `json:"id"`
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
