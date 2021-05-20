package endpoints

import (
	"coinswallet/internal/domain/wallet/entity"
)

type (
	//
	// CreateAccountRequest - holds the request params for the CreateAccount method
	CreateAccountRequest struct {
		Name entity.AccountName
	}
	// CreateAccountRequest - holds the response values for the CreateAccount method
	CreateAccountResponse struct {
		ID  entity.AccountID `json:"account_id,omitempty"`
		Err error            `json:"error,omitempty"`
	}

	//
	// DepositRequest - holds the request params for the Deposit method
	DepositRequest struct {
		Name   entity.AccountName
		Amount float64
	}
	// DepositResponse - holds the response values for the Deposit method
	DepositResponse struct {
		Balance float64 `json:"balance,omitempty"`
		Err     error   `json:"error,omitempty"`
	}

	//
	// TransferRequest - holds the request params for the Transfer method
	TransferRequest struct {
		From   entity.AccountName
		To     entity.AccountName
		Amount float64
	}
	// TransferResponse - holds the response values for the Transfer method
	TransferResponse struct {
		ID  entity.ID `json:"payment_id,omitempty"`
		Err error     `json:"error,omitempty"`
	}

	//
	// PaymentsListRequest - holds the request params for the PaymentsList method
	PaymentsListRequest struct {
		Name   entity.AccountName
		Offset int64
		Limit  int64
	}
	// PaymentsListResponse - holds the response values for the PaymentsList method
	PaymentsListResponse struct {
		List []entity.Payment `json:"list"`
		Err  error            `json:"error,omitempty"`
	}

	//
	// PaymentsListRequest - holds the request params for the AllPaymentsList method
	AllPaymentsListRequest struct {
		Offset int64
		Limit  int64
	}
	// PaymentsListResponse - holds the response values for the AllPaymentsList method
	AllPaymentsListResponse struct {
		List []entity.Payment `json:"list"`
		Err  error            `json:"error,omitempty"`
	}

	//
	// AccountsListRequest - holds the request params for the AccountsList method
	AccountsListRequest struct {
		Offset int64
		Limit  int64
	}
	// AccountsListResponse - holds the response values for the AccountsList method
	AccountsListResponse struct {
		List []entity.AccountName `json:"list"`
		Err  error                `json:"error,omitempty"`
	}
)
