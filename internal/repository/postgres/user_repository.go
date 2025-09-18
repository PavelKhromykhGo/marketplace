package postgres

import (
	"context"
	"marketplace/internal/user"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *user.User) error {
	_, err := r.db.NamedExecContext(ctx, `
INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
VALUES (:username, :email, :password_hash, :role, now(), now())
`, u)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*user.User, error) {
	var u user.User
	err := r.db.GetContext(ctx, &u, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	err := r.db.GetContext(ctx, &u, "SELECT * FROM users WHERE username=$1", username)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Update(ctx context.Context, user *user.User) error {
	_, err := r.db.NamedExecContext(ctx, `
UPDATE users SET username=:username, email=:email, password_hash=:password_hash, role=:role, updated_at=:updated_at
WHERE id=:id
`, user)
	return err
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	return err
}

func (r *UserRepo) UpdateRoleAndPassword(ctx context.Context, id int64, role, passhash string) error {
	_, err := r.db.ExecContext(ctx, `
UPDATE users SET role=$1, password_hash=$2, updated_at=now()
WHERE id=$3
`, role, passhash, id)
	return err
}
