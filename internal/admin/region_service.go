package admin

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/rashidmailru/kabobfood/internal/regions"
)

// RegionService manages region CRUD with cache invalidation.
type RegionService struct {
	repo  *regions.Repository
	cache *redis.Client
}

func NewRegionService(repo *regions.Repository, cache *redis.Client) *RegionService {
	return &RegionService{repo: repo, cache: cache}
}

func (s *RegionService) invalidate(ctx context.Context) {
	if s.cache == nil {
		return
	}
	_ = s.cache.Del(ctx, "regions:v1").Err()
	_ = s.cache.Del(ctx, "menu:v1").Err()
}

func (s *RegionService) CreateRegion(ctx context.Context, region regions.Region) (*regions.Region, error) {
	created, err := s.repo.Insert(ctx, region)
	if err != nil {
		return nil, err
	}
	s.invalidate(ctx)
	return created, nil
}

func (s *RegionService) UpdateRegion(ctx context.Context, region regions.Region) (*regions.Region, error) {
	updated, err := s.repo.Update(ctx, region)
	if err != nil {
		return nil, err
	}
	s.invalidate(ctx)
	return updated, nil
}

func (s *RegionService) DeleteRegion(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidate(ctx)
	return nil
}
