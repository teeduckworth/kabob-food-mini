package profile

import (
	"context"

	"github.com/rashidmailru/kabobfood/internal/addresses"
	"github.com/rashidmailru/kabobfood/internal/users"
)

// Service aggregates data for profile endpoints.
type Service struct {
	userRepo   *users.Repository
	addressSrv *addresses.Service
}

// NewService builds profile service.
func NewService(userRepo *users.Repository, addressService *addresses.Service) *Service {
	return &Service{userRepo: userRepo, addressSrv: addressService}
}

// Profile contains summary info for the current user.
type Profile struct {
	User      *users.User         `json:"user"`
	Addresses []addresses.Address `json:"addresses"`
}

// GetProfile returns profile info for a user.
func (s *Service) GetProfile(ctx context.Context, userID int64) (*Profile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	list, err := s.addressSrv.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &Profile{User: user, Addresses: list}, nil
}
