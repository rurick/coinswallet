package service

import (
	"coinswallet/internal/domain/wallet/entity"
	"context"
	"github.com/go-kit/kit/log"
	"os"
	"testing"
)

var logger log.Logger

func initLogger() {
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
}

func Test_CreateAccount(t *testing.T) {
	const validAccName = "Testing987ha9871hgaf98782"
	const invalidAccName = "Testing987ha9 871hgaf98782"
	initLogger()

	srv := newService(logger)
	t.Run("with valid account name", func(t *testing.T) {
		if err := srv.CreateAccount(context.Background(), validAccName); err != nil {
			t.Error(err)
		}
		// delete account
		a, err := entity.NewAccount()
		if err != nil {
			t.Fatal(err)
		}
		_ = a.Find(validAccName)
		_ = a.Delete()
	})
	t.Run("with invalid account name", func(t *testing.T) {
		if err := srv.CreateAccount(context.Background(), invalidAccName); err == nil {
			t.Error("wait error but not")
		}
	})
}

func Test_Deposit(t *testing.T) {
	const validAccName = "Testing987ha9871hgaf98782"
	initLogger()

	a, err := entity.NewAccount()
	if err != nil {
		t.Fatal(err)
	}
	srv := newService(logger)

	t.Run("create temp account", func(t *testing.T) {
		if err := a.Register(validAccName); err != nil {
			t.Error(err)
		}
	})
	t.Run("run service ", func(t *testing.T) {
		if err := srv.Deposit(context.Background(), validAccName, 2); err != nil {
			t.Error(err)
		}
	})
	t.Run("delete temp account", func(t *testing.T) {
		if err := a.Delete(); err != nil {
			t.Error(err)
		}
	})
}

func Test_Transfer(t *testing.T) {
	const validAccName1 = "Testing987ha9871hgaf987821"
	const validAccName2 = "Testing987ha9871hgaf987822"
	initLogger()

	a1, err := entity.NewAccount()
	if err != nil {
		t.Fatal(err)
	}
	a2, err := entity.NewAccount()
	srv := newService(logger)

	t.Run("create temp account", func(t *testing.T) {
		if err := a1.Register(validAccName1); err != nil {
			t.Error(err)
		}
	})
	t.Run("run service ", func(t *testing.T) {
		if err := srv.Deposit(context.Background(), validAccName, 2); err != nil {
			t.Error(err)
		}
	})
	t.Run("delete temp account", func(t *testing.T) {
		if err := a.Delete(); err != nil {
			t.Error(err)
		}
	})
}
