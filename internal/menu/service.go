package menu

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/rashidmailru/kabobfood/internal/regions"
)

const (
	menuCacheKey    = "menu:v1"
	regionsCacheKey = "regions:v1"
)

// Service aggregates menu and region data with caching.
type Service struct {
	menuRepo   MenuRepository
	regionRepo RegionRepository
	cache      *redis.Client
	menuTTL    time.Duration
	regionsTTL time.Duration
}

// ServiceConfig holds dependencies.
type ServiceConfig struct {
	MenuRepo   MenuRepository
	RegionRepo RegionRepository
	Cache      *redis.Client
	MenuTTL    time.Duration
	RegionsTTL time.Duration
}

// MenuRepository defines menu storage access methods.
type MenuRepository interface {
	GetActiveCategories(ctx context.Context) ([]Category, error)
	GetActiveProducts(ctx context.Context) ([]Product, error)
}

// RegionRepository defines region storage access methods.
type RegionRepository interface {
	GetActiveRegions(ctx context.Context) ([]regions.Region, error)
}

// NewService constructs service.
func NewService(cfg ServiceConfig) *Service {
	if cfg.MenuRepo == nil {
		panic("menu service: menu repository is required")
	}
	if cfg.RegionRepo == nil {
		panic("menu service: region repository is required")
	}
	ttlMenu := cfg.MenuTTL
	if ttlMenu <= 0 {
		ttlMenu = 30 * time.Second
	}
	ttlRegions := cfg.RegionsTTL
	if ttlRegions <= 0 {
		ttlRegions = 30 * time.Second
	}

	return &Service{
		menuRepo:   cfg.MenuRepo,
		regionRepo: cfg.RegionRepo,
		cache:      cfg.Cache,
		menuTTL:    ttlMenu,
		regionsTTL: ttlRegions,
	}
}

// GetMenu returns categories with products (cached).
func (s *Service) GetMenu(ctx context.Context) (*MenuResponse, error) {
	if s.cache != nil {
		if data, err := s.cache.Get(ctx, menuCacheKey).Bytes(); err == nil {
			var resp MenuResponse
			if err := json.Unmarshal(data, &resp); err == nil {
				return &resp, nil
			}
		}
	}

	categories, err := s.menuRepo.GetActiveCategories(ctx)
	if err != nil {
		return nil, err
	}
	products, err := s.menuRepo.GetActiveProducts(ctx)
	if err != nil {
		return nil, err
	}

	resp := buildMenuResponse(categories, products)

	if s.cache != nil {
		if bytes, err := json.Marshal(resp); err == nil {
			_ = s.cache.Set(ctx, menuCacheKey, bytes, s.menuTTL).Err()
		}
	}

	return resp, nil
}

// GetRegions returns active regions (cached).
func (s *Service) GetRegions(ctx context.Context) ([]regions.Region, error) {
	if s.cache != nil {
		if data, err := s.cache.Get(ctx, regionsCacheKey).Bytes(); err == nil {
			var resp []regions.Region
			if err := json.Unmarshal(data, &resp); err == nil {
				return resp, nil
			}
		}
	}

	list, err := s.regionRepo.GetActiveRegions(ctx)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		if bytes, err := json.Marshal(list); err == nil {
			_ = s.cache.Set(ctx, regionsCacheKey, bytes, s.regionsTTL).Err()
		}
	}

	return list, nil
}

func buildMenuResponse(categories []Category, products []Product) *MenuResponse {
	ordered := make([]MenuCategory, len(categories))
	indexByID := make(map[int64]int, len(categories))
	for i, cat := range categories {
		ordered[i] = MenuCategory{
			ID:        cat.ID,
			Name:      cat.Name,
			Emoji:     cat.Emoji,
			SortOrder: cat.SortOrder,
			Products:  []Product{},
		}
		indexByID[cat.ID] = i
	}

	for _, product := range products {
		if idx, ok := indexByID[product.CategoryID]; ok {
			ordered[idx].Products = append(ordered[idx].Products, product)
		}
	}

	return &MenuResponse{Categories: ordered}
}
