package store

import (
	"context"
	"errors"

	"github.com/Vialmsi/Interview/internal/entity"

	_ "github.com/lib/pq"
)

// NewUser function which allows to create (register) new user (register)
func (s *Store) NewUser(ctx context.Context, user entity.User) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	query := `INSERT INTO users(login,password,name,email) VALUES($1,$2,$3,$4) RETURNING ID`

	var id int

	err := s.db.QueryRowContext(ctx, query,
		user.Login,
		user.Password,
		user.Name,
		user.Email).Scan(&id)

	if isUniqueViolation(err) {
		return 0, errors.New("user with this login already exists")
	}

	return id, err
}

// RetrieveUser func which allows to get user using login and password (login)
func (s *Store) RetrieveUser(ctx context.Context, login, password string) (entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	query := `SELECT id,name,email FROM users WHERE login=$1 AND password=$2`

	var user entity.User

	err := s.db.QueryRowContext(ctx, query, login, password).Scan(
		&user.ID,
		&user.Name,
		&user.Email)

	return user, err
}
