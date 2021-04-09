package binarylane

import (
	"context"
	"net/http"
	"time"
)

// BalanceService is an interface for interfacing with the Balance
// endpoints of the BinaryLane API
// See: https://api.binarylane.com.au/reference/#balance
type BalanceService interface {
	Get(context.Context) (*Balance, *Response, error)
}

// BalanceServiceOp handles communication with the Balance related methods of
// the BinaryLane API.
type BalanceServiceOp struct {
	client *Client
}

var _ BalanceService = &BalanceServiceOp{}

// Balance represents a BinaryLane Balance
type Balance struct {
	MonthToDateBalance string    `json:"month_to_date_balance"`
	AccountBalance     string    `json:"account_balance"`
	MonthToDateUsage   string    `json:"month_to_date_usage"`
	GeneratedAt        time.Time `json:"generated_at"`
}

func (r Balance) String() string {
	return Stringify(r)
}

// Get balance info
func (s *BalanceServiceOp) Get(ctx context.Context) (*Balance, *Response, error) {
	path := "v2/customers/my/balance"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Balance)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
