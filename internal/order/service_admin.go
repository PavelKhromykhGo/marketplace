package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrOrderNotFound           = errors.New("order not found")
)

type AdminService interface {
	Ship(ctx context.Context, orderID int64) error
	Deliver(ctx context.Context, orderID int64) error
	Cancel(ctx context.Context, orderID int64) error
}

func (s *service) Ship(ctx context.Context, orderID int64) error {
	from, err := s.repoStatus(ctx, orderID)
	if err != nil {
		return err
	}
	to := StatusShipped
	if !IsValidStatusTransition(from, to) {
		return fmt.Errorf("%w: %s --> %s", ErrInvalidStatusTransition, from, to)
	}
	return s.repo.UpdateOrderStatus(ctx, orderID, from, to)
}

func (s *service) Deliver(ctx context.Context, orderID int64) error {
	from, err := s.repoStatus(ctx, orderID)
	if err != nil {
		return err
	}
	to := StatusDelivered
	if !IsValidStatusTransition(from, to) {
		return fmt.Errorf("%w: %s --> %s", ErrInvalidStatusTransition, from, to)
	}
	return s.repo.UpdateOrderStatus(ctx, orderID, from, to)
}

func (s *service) Cancel(ctx context.Context, orderID int64) error {
	from, err := s.repoStatus(ctx, orderID)
	if err != nil {
		return err
	}
	to := StatusCancelled
	if !IsValidStatusTransition(from, to) {
		return fmt.Errorf("%w: %s --> %s", ErrInvalidStatusTransition, from, to)
	}
	return s.repo.UpdateOrderStatus(ctx, orderID, from, to)
}

func (s *service) repoStatus(ctx context.Context, orderID int64) (string, error) {
	if getter, ok := s.repo.(interface {
		GetOrderStatus(context.Context, int64) (string, error)
	}); ok {
		st, err := getter.GetOrderStatus(ctx, orderID)
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrOrderNotFound
		}
		return st, err
	}
	return "", errors.New("repository does not support status reading")
}
