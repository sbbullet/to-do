package db

import (
	"database/sql"

	"github.com/google/uuid"
)

type CreateTodoParams struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Title    string    `json:"title"`
}

func (store *Store) CreateTodo(arg CreateTodoParams) (todo Todo, err error) {
	const createTodoQuery = `
		INSERT INTO todos(id, username, title)
		VALUES(?, ?, ?)
		RETURNING id, username, title, is_completed, created_at;
	`

	row := store.DB.QueryRow(createTodoQuery, arg.ID, arg.Username, arg.Title)

	err = row.Scan(&todo.ID, &todo.Username, &todo.Title, &todo.IsCompleted, &todo.CreatedAt)

	return
}

type GetUserTodosParams struct {
	Username string
	Limit    int
	Offset   int
}

func (store *Store) GetTodoById(id uuid.UUID) (todo Todo, err error) {
	const getTodoByIdQuery = `
		SELECT id, username, title, is_completed, created_at
		FROM todos
		WHERE id = ?;
	`

	row := store.DB.QueryRow(getTodoByIdQuery, id)

	err = row.Scan(&todo.ID, &todo.Username, &todo.Title, &todo.IsCompleted, &todo.CreatedAt)

	return
}

func (store *Store) GetUserTodos(arg GetUserTodosParams) ([]Todo, error) {
	const getUserTodosQuery = `
		SELECT id, username, title, is_completed, created_at
		FROM todos
		WHERE username = ?
		ORDER BY created_at
		LIMIT ?
		OFFSET ?;
	`

	rows, err := store.DB.Query(getUserTodosQuery, arg.Username, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var todo Todo

		err := rows.Scan(&todo.ID, &todo.Username, &todo.Title, &todo.IsCompleted, &todo.CreatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

type UpdateTodoParams struct {
	ID          uuid.UUID      `json:"id"`
	Title       sql.NullString `json:"title"`
	IsCompleted sql.NullBool   `json:"is_completed"`
}

func (store *Store) UpdateTodo(arg UpdateTodoParams) (todo Todo, err error) {
	const updateTodoQuery = `
		UPDATE todos
		SET
			title = COALESCE(?, title),
			is_completed = COALESCE(?, is_completed)
		WHERE
			id = ?
		RETURNING id, username, title, is_completed, created_at;
	`

	row := store.DB.QueryRow(updateTodoQuery, arg.Title, arg.IsCompleted, arg.ID)

	err = row.Scan(&todo.ID, &todo.Username, &todo.Title, &todo.IsCompleted, &todo.CreatedAt)

	return
}

type DeleteTodoOfAUserParams struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func (store *Store) DeleteTodoOfAUser(arg DeleteTodoOfAUserParams) error {
	const deleteTodoByIdQuery = `
		DELETE FROM todos
		WHERE id = ? AND username = ?;
	`

	result, err := store.DB.Exec(deleteTodoByIdQuery, arg.ID, arg.Username)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		return sql.ErrNoRows
	}

	return nil
}
