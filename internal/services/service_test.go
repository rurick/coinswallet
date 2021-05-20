// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// yuriyiovkov@gmail.com; telegram: @yuriyiovkov

// Important!
// for successfully test passed run it from directory where file .env is
// or set up the ENV in your OS environment

package services

import (
	"context"
	"os"
	"testing"

	"coinswallet/internal/domain/wallet/entity"
	"github.com/go-kit/kit/log"
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

	srv := NewService(logger)
	t.Run("with valid account name", func(t *testing.T) {
		if _, err := srv.CreateAccount(context.Background(), validAccName); err != nil {
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
		if _, err := srv.CreateAccount(context.Background(), invalidAccName); err == nil {
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
	srv := NewService(logger)

	t.Run("create temp account", func(t *testing.T) {
		if err := a.Register(validAccName); err != nil {
			t.Error(err)
		}
	})
	t.Run("run service ", func(t *testing.T) {
		if _, err := srv.Deposit(context.Background(), validAccName, 3); err != nil {
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
	srv := NewService(logger)

	t.Run("create temp account 1", func(t *testing.T) {
		if err := a1.Register(validAccName1); err != nil {
			t.Error(err)
		}
	})
	t.Run("create temp account 2", func(t *testing.T) {
		if err := a2.Register(validAccName2); err != nil {
			t.Error(err)
		}
	})
	t.Run("deposit ", func(t *testing.T) {
		if _, err := srv.Deposit(context.Background(), validAccName1, 2); err != nil {
			t.Error(err)
		}
	})
	t.Run("run service ", func(t *testing.T) {
		if _, err := srv.Transfer(context.Background(), validAccName1, validAccName2, 1); err != nil {
			t.Error(err)
		}
		_ = a1.Find(validAccName1)
		_ = a2.Find(validAccName2)
		if a1.Balance != a2.Balance && a2.Balance != 1 {
			t.Errorf("balanse must be 1 and 1, have: %f, %f", a1.Balance, a2.Balance)
		}
	})

	t.Run("delete temp accounts", func(t *testing.T) {
		if err := a1.Delete(); err != nil {
			t.Error(err)
		}
		if err := a2.Delete(); err != nil {
			t.Error(err)
		}
	})
}

func Test_PaymentsList(t *testing.T) {
	const validAccName = "Testing987ha9871hgaf92c8782"
	initLogger()

	a, err := entity.NewAccount()
	if err != nil {
		t.Fatal(err)
	}
	srv := NewService(logger)

	t.Run("create temp account", func(t *testing.T) {
		if err := a.Register(validAccName); err != nil {
			t.Error(err)
		}
	})
	t.Run("deposit ", func(t *testing.T) {
		if _, err := srv.Deposit(context.Background(), validAccName, 6); err != nil {
			t.Error(err)
		}
	})
	t.Run("run service ", func(t *testing.T) {
		var lst []PaymentEntity
		if lst, err = srv.PaymentsList(context.Background(), validAccName, 0, -1); err != nil {
			t.Error(err)
		}
		if len(lst) != 1 {
			t.Errorf("list length must be 1, got: %d", len(lst))
		}
	})

	t.Run("delete temp account", func(t *testing.T) {
		if err := a.Delete(); err != nil {
			t.Error(err)
		}
	})
}

func Test_AllPaymentsList(t *testing.T) {
	initLogger()

	srv := NewService(logger)

	t.Run("run service ", func(t *testing.T) {
		if _, err := srv.AllPaymentsList(context.Background(), 0, -1); err != nil {
			t.Error(err)
		}
	})
}

func Test_AccountsList(t *testing.T) {
	initLogger()

	srv := NewService(logger)

	t.Run("run service ", func(t *testing.T) {
		if _, err := srv.AccountsList(context.Background(), 0, -1); err != nil {
			t.Error(err)
		}
	})
}
