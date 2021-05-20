package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"coinswallet/internal/domain/wallet/entity"
	"coinswallet/internal/endpoints"
	"coinswallet/internal/services"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPHandler(s services.Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := endpoints.MakeEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST 	/account/						create new wallet account
	// PATCH 	/account/deposit/				deposit amount of currency to the wallet account
	// PATCH 	/account/transfer/				send amount of currency between two wallet accounts
	// GET	 	/accounts/:offset/:limit/		list of all registered accounts
	// GET	 	/payments/:name:/offset/:limit/	list of payments of the account
	// GET	 	/payments/:offset/:limit/		list of all payments

	r.Methods("POST").Path("/account/").Handler(httptransport.NewServer(
		e.CreateAccount,
		decodeCreateAccount,
		encodeResponse,
		options...,
	))
	r.Methods("PATCH").Path("/account/deposit/").Handler(httptransport.NewServer(
		e.Deposit,
		decodeDeposit,
		encodeResponse,
		options...,
	))
	r.Methods("PATCH").Path("/account/transfer/").Handler(httptransport.NewServer(
		e.Transfer,
		decodeTransfer,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/payments/{name}/{offset}/{limit}/").Handler(httptransport.NewServer(
		e.PaymentsList,
		decodePaymentsList,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/payments/{offset}/{limit}/").Handler(httptransport.NewServer(
		e.AllPaymentsList,
		decodeAllPaymentsList,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/accounts/{offset}/{limit}/").Handler(httptransport.NewServer(
		e.AccountsList,
		decodeAccountsList,
		encodeResponse,
		options...,
	))
	return r
}

func decodeCreateAccount(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req endpoints.CreateAccountRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeDeposit(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req endpoints.DepositRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeTransfer(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req endpoints.TransferRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodePaymentsList(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, ErrBadRouting
	}

	o, ok := vars["offset"]
	if !ok {
		return nil, ErrBadRouting
	}
	offset, e := strconv.ParseInt(o, 10, 64)
	if e != nil {
		return nil, ErrBadRouting
	}

	l, ok := vars["limit"]
	if !ok {
		return nil, ErrBadRouting
	}
	limit, e := strconv.ParseInt(l, 10, 64)
	if e != nil {
		return nil, ErrBadRouting
	}

	return endpoints.PaymentsListRequest{Name: entity.AccountName(name), Offset: offset, Limit: limit}, nil
}

func decodeAllPaymentsList(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	o, ok := vars["offset"]
	if !ok {
		return nil, ErrBadRouting
	}
	offset, e := strconv.ParseInt(o, 10, 64)
	if e != nil {
		return nil, ErrBadRouting
	}

	l, ok := vars["limit"]
	if !ok {
		return nil, ErrBadRouting
	}
	limit, e := strconv.ParseInt(l, 10, 64)
	if e != nil {
		return nil, ErrBadRouting
	}

	return endpoints.AllPaymentsListRequest{Offset: offset, Limit: limit}, nil
}

func decodeAccountsList(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	o, ok := vars["offset"]
	if !ok {
		return nil, ErrBadRouting
	}
	offset, e := strconv.ParseInt(o, 10, 64)
	if e != nil {
		return nil, ErrBadRouting
	}

	l, ok := vars["limit"]
	if !ok {
		return nil, ErrBadRouting
	}
	limit, e := strconv.ParseInt(l, 10, 64)
	if e != nil {
		return nil, ErrBadRouting
	}

	return endpoints.AccountsListRequest{Offset: offset, Limit: limit}, nil
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.Error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	Error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case services.ErrDepositNotFound,
		services.ErrTransferFromNotFound,
		services.ErrTransferToNotFound,
		services.ErrPaymentsListNotFound:

		return http.StatusNotFound

	case services.ErrCreateAccountInvalidName,
		services.ErrCreateAccount,
		services.ErrDepositAmountError,
		services.ErrTransferAmountError,
		services.ErrPaymentsListOffsetLimitError,
		services.ErrAccountsListOffsetLimitError:

		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}
