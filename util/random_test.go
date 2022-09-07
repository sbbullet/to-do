package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	const (
		min int64 = 5
		max int64 = 30
	)

	value := RandomInt(min, max)

	require.GreaterOrEqual(t, value, min)
	require.LessOrEqual(t, value, max)
}

func TestRandomString(t *testing.T) {
	length := 8
	randomString := RandomString(length)
	require.Equal(t, len(randomString), length)
}
