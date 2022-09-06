package db

import "time"

type User struct {
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	FullName       string    `json:"full_name"`
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
}

type Todo struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Title       string    `json:"title"`
	CreatedAt   time.Time `json:"created_at"`
	IsCompleted bool      `json:"is_completed"`
}
