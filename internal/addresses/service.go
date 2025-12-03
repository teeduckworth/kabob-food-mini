package addresses

import "context"

// Service contains business logic around addresses.
type Service struct {
	repo *Repository
}

// NewService creates service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List returns all addresses for user.
func (s *Service) List(ctx context.Context, userID int64) ([]Address, error) {
	return s.repo.ListByUser(ctx, userID)
}

// Create adds a new address and manages default flag.
func (s *Service) Create(ctx context.Context, input CreateInput) (*Address, error) {
	if input.IsDefault {
		if err := s.repo.ClearDefault(ctx, input.UserID); err != nil {
			return nil, err
		}
	}
	return s.repo.Insert(ctx, input)
}

// Update modifies address fields and default flag.
func (s *Service) Update(ctx context.Context, input UpdateInput) (*Address, error) {
	if input.IsDefault {
		if err := s.repo.ClearDefault(ctx, input.UserID); err != nil {
			return nil, err
		}
	}
	return s.repo.Update(ctx, input)
}

// Delete removes address.
func (s *Service) Delete(ctx context.Context, id, userID int64) error {
	return s.repo.Delete(ctx, id, userID)
}
