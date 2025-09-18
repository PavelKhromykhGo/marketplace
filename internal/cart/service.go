package cart

import "context"

type Repository interface {
	AddItem(ctx context.Context, item *CartItem) (int64, error)
	ListItems(ctx context.Context, userID int64) ([]*CartItem, error)
	RemoveItem(ctx context.Context, userID, productID int64) error
	Clear(ctx context.Context, userID int64) error
}

type Service interface {
	AddItem(ctx context.Context, userID, productID int64, qty int) (int64, error)
	ListItem(ctx context.Context, userID int64) ([]*CartItem, error)
	RemoveItem(ctx context.Context, userID, productID int64) error
	Clear(ctx context.Context, userID int64) error
}

type cartService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &cartService{repo: repo}
}

func (c *cartService) AddItem(ctx context.Context, userID, productID int64, qty int) (int64, error) {
	item := &CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  int64(qty),
	}
	return c.repo.AddItem(ctx, item)
}

func (c *cartService) ListItem(ctx context.Context, userID int64) ([]*CartItem, error) {
	return c.repo.ListItems(ctx, userID)
}

func (c *cartService) RemoveItem(ctx context.Context, userID, productID int64) error {
	return c.repo.RemoveItem(ctx, userID, productID)
}

func (c *cartService) Clear(ctx context.Context, userID int64) error {
	return c.repo.Clear(ctx, userID)
}
