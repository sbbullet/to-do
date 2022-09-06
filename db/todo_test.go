package db

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sbbullet/to-do/util"
	"github.com/stretchr/testify/require"
)

func TestCreateTodo(t *testing.T) {
	user := createRandomUser(t)
	createRandomTodo(t, user.Username)
}

func TestGetUserTodos(t *testing.T) {
	user := createRandomUser(t)

	n := 5
	var todo Todo
	for i := 1; i <= n; i++ {
		todo = createRandomTodo(t, user.Username)
	}

	userTodos, err := testStore.GetUserTodos(GetUserTodosParams{
		Username: user.Username,
		Limit:    n,
		Offset:   0,
	})

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(userTodos), n)
	require.Equal(t, todo.Username, user.Username)
	require.Contains(t, userTodos, todo)
}

func createRandomTodo(t *testing.T, username string) Todo {
	todoID, err := uuid.NewRandom()
	require.NoError(t, err)
	require.NotEmpty(t, todoID)

	arg := CreateTodoParams{
		ID:       todoID.String(),
		Username: username,
		Title:    util.RandomString(50),
	}

	todo, err := testStore.CreateTodo(arg)
	require.NoError(t, err)
	require.NotEmpty(t, todo)

	require.Equal(t, arg.ID, todo.ID)
	require.Equal(t, arg.Username, todo.Username)
	require.Equal(t, arg.Title, todo.Title)
	require.False(t, todo.IsCompleted)
	require.WithinDuration(t, time.Now(), todo.CreatedAt, 2*time.Second)

	return todo
}
