package payment

import (
	"context"
	"marketplace/internal/order"
)

type Repository interface {
	CreateIntent(ctx context.Context, o *order.Order) (*Intent, error)
	ConfirmIntent(ctx context.Context, orderID int64, clientSecret string) (*Intent, error)
}

type OrderRepository interface {
	GetOrderWithItems(ctx context.Context, userID, orderID int64) (*order.Order, error)
}

type Service struct {
	payRepo Repository
	ordRepo OrderRepository
}

func NewService(payRepo Repository, ordRepo OrderRepository) *Service {
	return &Service{payRepo: payRepo, ordRepo: ordRepo}
}

func (s *Service) CreateIntent(ctx context.Context, userID, orderID int64) (*Intent, error) {
	o, err := s.ordRepo.GetOrderWithItems(ctx, userID, orderID)
	if err != nil {
		return nil, err
	}
	return s.payRepo.CreateIntent(ctx, o)
}

func (s *Service) Confirm(ctx context.Context, userID, orderID int64, clientSecret string) (*Intent, error) {
	if _, err := s.ordRepo.GetOrderWithItems(ctx, userID, orderID); err != nil {
		return nil, err
	}
	return s.payRepo.ConfirmIntent(ctx, orderID, clientSecret)
}
