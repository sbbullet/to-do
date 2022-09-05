package db

type CreateUserParams struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	FullName       string `json:"full_name"`
	HashedPassword string `json:"hashed_password"`
}

func (store *Store) CreateUser(arg CreateUserParams) (user User, err error) {
	const createUserQuery = `
		INSERT INTO users(username, email, full_name, hashed_password)
		VALUES(?, ?, ?, ?)
		RETURNING *;
	`

	row := store.DB.QueryRow(createUserQuery,
		arg.Username,
		arg.Email,
		arg.FullName,
		arg.HashedPassword,
	)

	err = row.Scan(
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.HashedPassword,
		&user.CreatedAt,
	)

	return
}

func (store *Store) GetUser(username string) (user User, err error) {
	const getUserQuery = `
		SELECT * FROM users
		WHERE username = ?;
		`

	row := store.DB.QueryRow(getUserQuery, username)

	err = row.Scan(&user.Username, &user.Email, &user.FullName, &user.HashedPassword, &user.CreatedAt)

	return
}
