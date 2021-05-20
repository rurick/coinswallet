package endpoints

import (
	"context"

	"coinswallet/internal/service"
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
func MakeEndpoints(s service.Service) Endpoints {
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

func makeCreateAccountEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateAccountRequest)
		err := s.CreateAccount(ctx, req.Name)
		return CreateAccountResponse{Err: err}, nil
	}
}
func makeDepositEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DepositRequest)
		err := s.Deposit(ctx, req.Name, req.Amount)
		return DepositResponse{Err: err}, nil
	}
}
func makeTransferEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TransferRequest)
		err := s.Transfer(ctx, req.From, req.To, req.Amount)
		return TransferResponse{Err: err}, nil
	}
}
func makePaymentsListEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PaymentsListRequest)
		lst, err := s.PaymentsList(ctx, req.Name, req.Offset, req.Limit)
		return PaymentsListResponse{List: lst, Err: err}, nil
	}
}
func makeAllPaymentsListEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AllPaymentsListRequest)
		lst, err := s.AllPaymentsList(ctx, req.Offset, req.Limit)
		return AllPaymentsListResponse{List: lst, Err: err}, nil
	}
}
func makeAccountsListEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AccountsListRequest)
		lst, err := s.AccountsList(ctx, req.Offset, req.Limit)
		return AccountsListResponse{List: lst, Err: err}, nil
	}
}
