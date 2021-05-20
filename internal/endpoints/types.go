package endpoints

import (
	"coinswallet/internal/domain/wallet/entity"
)

//
// CreateAccountRequest - holds the request params for the CreateAccount method
type CreateAccountRequest struct {
	Name entity.AccountName
}

// CreateAccountResponse - holds the response values for the CreateAccount method
type CreateAccountResponse struct {
	ID  entity.AccountName `json:"account_id,omitempty"`
	Err error              `json:"error,omitempty"`
}

func (r CreateAccountResponse) Error() error { return r.Err }

//
// DepositRequest - holds the request params for the Deposit method
type DepositRequest struct {
	Name   entity.AccountName
	Amount float64
}

// DepositResponse - holds the response values for the Deposit method
type DepositResponse struct {
	Balance float64 `json:"balance,omitempty"`
	Err     error   `json:"error,omitempty"`
}

func (r DepositResponse) Error() error { return r.Err }

//
// TransferRequest - holds the request params for the Transfer method
type TransferRequest struct {
	From   entity.AccountName
	To     entity.AccountName
	Amount float64
}

// TransferResponse - holds the response values for the Transfer method
type TransferResponse struct {
	Payment interface{} `json:"payment,omitempty"`
	Err     error       `json:"error,omitempty"`
}

func (r TransferResponse) Error() error { return r.Err }

//
// PaymentsListRequest - holds the request params for the PaymentsList method
type PaymentsListRequest struct {
	Name   entity.AccountName
	Offset int64
	Limit  int64
}

// PaymentsListResponse - holds the response values for the PaymentsList method
type PaymentsListResponse struct {
	List interface{} `json:"list"`
	Err  error       `json:"error,omitempty"`
}

func (r PaymentsListResponse) Error() error { return r.Err }

//
// AllPaymentsListRequest  - holds the request params for the AllPaymentsList method
type AllPaymentsListRequest struct {
	Offset int64
	Limit  int64
}

// AllPaymentsListResponse PaymentsListResponse - holds the response values for the AllPaymentsList method
type AllPaymentsListResponse struct {
	List interface{} `json:"list"`
	Err  error       `json:"error,omitempty"`
}

func (r AllPaymentsListResponse) Error() error { return r.Err }

//
// AccountsListRequest - holds the request params for the AccountsList method
type AccountsListRequest struct {
	Offset int64
	Limit  int64
}

// AccountsListResponse - holds the response values for the AccountsList method
type AccountsListResponse struct {
	List interface{} `json:"list"`
	Err  error       `json:"error,omitempty"`
}

func (r AccountsListResponse) Error() error { return r.Err }
