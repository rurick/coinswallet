package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"coinswallet/internal/endpoints"
	"coinswallet/internal/service"
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

func MakeHTTPHandler(s service.Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := endpoints.MakeEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST 	/account/						create new wallet account
	// PATCH 	/account/:name/deposit/			deposit amount of currency to the wallet account
	// PATCH 	/account/:name/transfer/		send amount of currency between two wallet accounts
	// GET	 	/accounts/						list of all registered accounts
	// GET	 	/payments/:name					list of payments of the account
	// GET	 	/payments/						list of all payments
	return r
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

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
	case service.ErrDepositNotFound,
		service.ErrTransferFromNotFound,
		service.ErrTransferToNotFound,
		service.ErrPaymentsListNotFound:

		return http.StatusNotFound

	case service.ErrCreateAccountInvalidName,
		service.ErrCreateAccount,
		service.ErrDepositAmountError,
		service.ErrTransferAmountError,
		service.ErrPaymentsListOffsetLimitError,
		service.ErrAccountsListOffsetLimitError:

		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}
