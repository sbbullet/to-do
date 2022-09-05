package db

import (
	"testing"
	"time"

	"github.com/sbbullet/to-do/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	createdUser := createRandomUser(t)

	resultingUser, err := testStore.GetUser(createdUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, resultingUser)

	require.Equal(t, resultingUser.Username, createdUser.Username)
	require.Equal(t, resultingUser.Email, createdUser.Email)
	require.Equal(t, resultingUser.FullName, createdUser.FullName)
	require.Equal(t, resultingUser.HashedPassword, createdUser.HashedPassword)
	require.WithinDuration(t, resultingUser.CreatedAt, createdUser.CreatedAt, time.Second)
}

// Create a random user in the test database
func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(8))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	arg := CreateUserParams{
		Username:       util.RandomUsername(),
		Email:          util.RandomEmail(),
		FullName:       util.RandomString(4) + " " + util.RandomString(4),
		HashedPassword: hashedPassword,
	}

	user, err := testStore.CreateUser(arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)

	return user
}
