package orders

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/rashidmailru/kabobfood/internal/addresses"
	"github.com/rashidmailru/kabobfood/internal/metrics"
	"github.com/rashidmailru/kabobfood/internal/notifications"
	"github.com/rashidmailru/kabobfood/internal/products"
	"github.com/rashidmailru/kabobfood/internal/regions"
	"github.com/rashidmailru/kabobfood/internal/users"
)

// Service handles order workflows.
type Service struct {
	repo        *Repository
	productRepo *products.Repository
	addressRepo *addresses.Repository
	regionRepo  *regions.Repository
	userRepo    *users.Repository
	notifier    *notifications.TelegramNotifier
	metrics     *metrics.Metrics
}

// NewService builds Service.
func NewService(repo *Repository, productRepo *products.Repository, addressRepo *addresses.Repository, regionRepo *regions.Repository, userRepo *users.Repository, notifier *notifications.TelegramNotifier, m *metrics.Metrics) *Service {
	return &Service{repo: repo, productRepo: productRepo, addressRepo: addressRepo, regionRepo: regionRepo, userRepo: userRepo, notifier: notifier, metrics: m}
}

var (
	errEmptyItems      = errors.New("order must contain at least one item")
	errInvalidType     = errors.New("order type must be delivery or pickup")
	errInvalidPayment  = errors.New("payment method is required")
	errInvalidRegion   = errors.New("region is not available")
	errInvalidAddress  = errors.New("address does not belong to user")
	errProductNotFound = errors.New("product not found or inactive")
	errInvalidQuantity = errors.New("item quantity must be > 0")
)

var allowedTypes = map[string]struct{}{"delivery": {}, "pickup": {}}

// Create creates an order ensuring idempotency and price recalculation.
func (s *Service) Create(ctx context.Context, userID int64, input CreateOrderInput) (*Order, error) {
	if _, err := uuid.Parse(input.ClientRequestID); err != nil {
		return nil, fmt.Errorf("invalid client_request_id: %w", err)
	}
	if len(input.Items) == 0 {
		return nil, errEmptyItems
	}
	orderType := strings.ToLower(input.Type)
	if _, ok := allowedTypes[orderType]; !ok {
		return nil, errInvalidType
	}
	if strings.TrimSpace(input.PaymentMethod) == "" {
		return nil, errInvalidPayment
	}

	if orderType == "delivery" && input.AddressID == 0 {
		return nil, errors.New("address is required for delivery")
	}

	region, err := s.regionRepo.GetByID(ctx, input.RegionID)
	if err != nil || !region.IsActive {
		return nil, errInvalidRegion
	}

	var addressID *int64
	if input.AddressID > 0 {
		addr, err := s.addressRepo.GetByIDAndUser(ctx, input.AddressID, userID)
		if err != nil {
			return nil, errInvalidAddress
		}
		addressID = &addr.ID
	}

	mergedItems := make(map[int64]int32)
	for _, item := range input.Items {
		if item.Qty <= 0 {
			return nil, errInvalidQuantity
		}
		mergedItems[item.ProductID] += item.Qty
	}

	productIDs := make([]int64, 0, len(mergedItems))
	for id := range mergedItems {
		productIDs = append(productIDs, id)
	}

	dbProducts, err := s.productRepo.GetActiveByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	orderItems := make([]OrderItem, 0, len(mergedItems))
	var itemsTotal float64
	for productID, qty := range mergedItems {
		product, ok := dbProducts[productID]
		if !ok || !product.IsActive {
			return nil, errProductNotFound
		}
		total := product.Price * float64(qty)
		itemsTotal += total
		orderItems = append(orderItems, OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Qty:         qty,
			Price:       product.Price,
			Total:       total,
		})
	}

	deliveryPrice := 0.0
	if orderType == "delivery" {
		deliveryPrice = region.DeliveryPrice
	}

	totalPrice := itemsTotal + deliveryPrice

	params := CreateParams{
		ClientRequestID: input.ClientRequestID,
		UserID:          userID,
		AddressID:       addressID,
		Type:            orderType,
		PaymentMethod:   input.PaymentMethod,
		Status:          "new",
		RegionID:        region.ID,
		DeliveryPrice:   deliveryPrice,
		ItemsTotal:      itemsTotal,
		TotalPrice:      totalPrice,
		Comment:         input.Comment,
		CustomerName:    input.CustomerName,
		CustomerPhone:   input.CustomerPhone,
		Items:           orderItems,
	}

	order, err := s.repo.Create(ctx, params)
	if err != nil {
		return nil, err
	}
	if s.metrics != nil {
		s.metrics.OrdersCreated.Inc()
	}
	if s.notifier != nil && s.userRepo != nil {
		if user, err := s.userRepo.GetByID(ctx, userID); err == nil && user.TelegramID != 0 {
			info := notifications.OrderInfo{
				OrderID:       order.ID,
				Status:        order.Status,
				Total:         order.TotalPrice,
				CustomerName:  order.CustomerName,
				CustomerPhone: order.CustomerPhone,
			}
			s.notifier.NotifyOrderCreated(ctx, info, user.TelegramID)
		}
	}

	return order, nil
}

// List returns user's orders.
func (s *Service) List(ctx context.Context, userID int64) ([]Order, error) {
	return s.repo.ListByUser(ctx, userID, 50)
}

// Get returns single order by id.
func (s *Service) Get(ctx context.Context, userID, orderID int64) (*Order, error) {
	return s.repo.GetByID(ctx, orderID, userID)
}

// AdminService exposes operations for operators.
type AdminService struct {
	repo     *Repository
	userRepo *users.Repository
	notifier *notifications.TelegramNotifier
}

// NewAdminService builds admin service.
func NewAdminService(repo *Repository, userRepo *users.Repository, notifier *notifications.TelegramNotifier) *AdminService {
	return &AdminService{repo: repo, userRepo: userRepo, notifier: notifier}
}

// List returns latest orders regardless of user.
func (s *AdminService) List(ctx context.Context, params AdminListParams) ([]Order, error) {
	return s.repo.ListAdmin(ctx, params)
}

// UpdateStatus updates and returns order.
func (s *AdminService) UpdateStatus(ctx context.Context, orderID int64, status string) (*Order, error) {
	allowed := map[string]struct{}{
		"new":       {},
		"accepted":  {},
		"cooking":   {},
		"delivery":  {},
		"delivered": {},
		"canceled":  {},
	}
	if _, ok := allowed[strings.ToLower(status)]; !ok {
		return nil, errors.New("invalid status")
	}
	order, err := s.repo.UpdateStatus(ctx, orderID, status)
	if err != nil {
		return nil, err
	}
	if s.notifier != nil && s.userRepo != nil {
		if user, err := s.userRepo.GetByID(ctx, order.UserID); err == nil && user.TelegramID != 0 {
			info := notifications.OrderInfo{
				OrderID:       order.ID,
				Status:        order.Status,
				Total:         order.TotalPrice,
				CustomerName:  order.CustomerName,
				CustomerPhone: order.CustomerPhone,
			}
			s.notifier.NotifyStatusChanged(ctx, info, user.TelegramID)
		}
	}
	return order, nil
}
