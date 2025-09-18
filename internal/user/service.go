package user

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id int64) error
	UpdateRoleAndPassword(ctx context.Context, id int64, role, passhash string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u := &User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         "user",
	}
	return s.repo.Create(ctx, u)
}

func (s *Service) Authenticate(ctx context.Context, username, password string) (*User, error) {
	u, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return u, nil
}

func (s *Service) EnsureAdmin(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return nil
	}

	u, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if u == nil || u.ID == 0 {
		return s.repo.Create(ctx, &User{
			Username:     username,
			PasswordHash: string(hash),
			Role:         "admin",
		})
	}
	if u.Role == "admin" && bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil {
		return nil
	}
	return s.repo.UpdateRoleAndPassword(ctx, u.ID, "admin", string(hash))
}
