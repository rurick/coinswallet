package endpoints

import (
	"context"

	"coinswallet/internal/services"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateAccount   endpoint.Endpoint
	Deposit         endpoint.Endpoint
	Transfer        endpoint.Endpoint
	PaymentsList    endpoint.Endpoint
	AllPaymentsList endpoint.Endpoint
	AccountsList    endpoint.Endpoint
}

// Endpoints holds all Go kit endpoints for the wallet service.
func MakeEndpoints(s services.Service) Endpoints {
	return Endpoints{
		CreateAccount:   makeCreateAccountEndpoint(s),
		Deposit:         makeDepositEndpoint(s),
		Transfer:        makeTransferEndpoint(s),
		PaymentsList:    makePaymentsListEndpoint(s),
		AllPaymentsList: makeAllPaymentsListEndpoint(s),
		AccountsList:    makeAccountsListEndpoint(s),
	}
}

//
// MakeEndpoints initializes all Go kit endpoints for the wallet service.
//

func makeCreateAccountEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateAccountRequest)
		id, err := s.CreateAccount(ctx, req.Name)
		return CreateAccountResponse{ID: id, Err: err}, nil
	}
}

func makeDepositEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DepositRequest)
		a, err := s.Deposit(ctx, req.Name, req.Amount)
		return DepositResponse{Balance: a, Err: err}, nil
	}
}

func makeTransferEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TransferRequest)
		id, err := s.Transfer(ctx, req.From, req.To, req.Amount)
		return TransferResponse{ID: id, Err: err}, nil
	}
}

func makePaymentsListEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PaymentsListRequest)
		lst, err := s.PaymentsList(ctx, req.Name, req.Offset, req.Limit)
		return PaymentsListResponse{List: lst, Err: err}, nil
	}
}

func makeAllPaymentsListEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AllPaymentsListRequest)
		lst, err := s.AllPaymentsList(ctx, req.Offset, req.Limit)
		return AllPaymentsListResponse{List: lst, Err: err}, nil
	}
}

func makeAccountsListEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AccountsListRequest)
		lst, err := s.AccountsList(ctx, req.Offset, req.Limit)
		return AccountsListResponse{List: lst, Err: err}, nil
	}
}
