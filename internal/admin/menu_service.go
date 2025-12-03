package admin

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/rashidmailru/kabobfood/internal/menu"
)

// MenuService wraps menu repo with cache invalidation.
type MenuService struct {
	repo  *menu.Repository
	cache *redis.Client
}

// NewMenuService builds admin menu service.
func NewMenuService(repo *menu.Repository, cache *redis.Client) *MenuService {
	return &MenuService{repo: repo, cache: cache}
}

func (s *MenuService) invalidateCache(ctx context.Context) {
	if s.cache == nil {
		return
	}
	keys := []string{"menu:v1"}
	for _, key := range keys {
		_ = s.cache.Del(ctx, key).Err()
	}
}

// CreateCategory creates a category and invalidates cache.
func (s *MenuService) CreateCategory(ctx context.Context, cat menu.Category) (*menu.Category, error) {
	created, err := s.repo.InsertCategory(ctx, cat)
	if err != nil {
		return nil, err
	}
	s.invalidateCache(ctx)
	return created, nil
}

// UpdateCategory updates category.
func (s *MenuService) UpdateCategory(ctx context.Context, cat menu.Category) (*menu.Category, error) {
	updated, err := s.repo.UpdateCategory(ctx, cat)
	if err != nil {
		return nil, err
	}
	s.invalidateCache(ctx)
	return updated, nil
}

// DeleteCategory removes category.
func (s *MenuService) DeleteCategory(ctx context.Context, id int64) error {
	if err := s.repo.DeleteCategory(ctx, id); err != nil {
		return err
	}
	s.invalidateCache(ctx)
	return nil
}

// CreateProduct creates product.
func (s *MenuService) CreateProduct(ctx context.Context, product menu.Product) (*menu.Product, error) {
	created, err := s.repo.InsertProduct(ctx, product)
	if err != nil {
		return nil, err
	}
	s.invalidateCache(ctx)
	return created, nil
}

// UpdateProduct updates product values.
func (s *MenuService) UpdateProduct(ctx context.Context, product menu.Product) (*menu.Product, error) {
	updated, err := s.repo.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}
	s.invalidateCache(ctx)
	return updated, nil
}

// DeleteProduct removes product.
func (s *MenuService) DeleteProduct(ctx context.Context, id int64) error {
	if err := s.repo.DeleteProduct(ctx, id); err != nil {
		return err
	}
	s.invalidateCache(ctx)
	return nil
}
