package db

import "fmt"

type CreateTodoParams struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Title    string `json:"title"`
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

func (store *Store) GetUserTodos(arg GetUserTodosParams) ([]Todo, error) {
	const getUserTodosQuery = `
		SELECT id, username, title, is_completed, created_at
		FROM todos
		WHERE username = ?
		ORDER BY created_at
		LIMIT ?
		OFFSET ?;
	`

	fmt.Println(arg)

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
