// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// yuriyiovkov@gmail.com; telegram: @yuriyiovkov

// Important!
// for successfully test passed run it from directory where file .env is
// or set up the ENV in your OS environment

package entity

import (
	"testing"
)

func Test_Create(t *testing.T) {
	type (
		args struct {
			name AccountName
			id   AccountID
		}
		dataSet []struct {
			name    string
			args    args
			wantErr bool
		}
	)
	tests := dataSet{
		{
			"with valid name",
			args{
				"mywallet_76ck76wecoan0vl",
				-1,
			},
			false,
		},
		{
			"with to long name",
			args{
				"asf4566asdkjhsakhkiwhckjashckakjcsdichsidcik",
				-1,
			},
			true,
		},
		{
			"with exists name",
			args{
				"mywallet_76ck76wecoan0vl",
				-1,
			},
			true,
		},
	}
	// create accounts
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewAccount()
			if err != nil {
				t.Fatal("NewAccount() error: ", err)
			}
			err = a.Register(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error : %v, wantErr %v ", err, tt.wantErr)
				return
			}
			tests[i].args.id = a.ID
		})
	}

	// delete created accounts
	t.Run("delete existing accounts", func(t *testing.T) {
		a, err := NewAccount()
		if err != nil {
			t.Fatalf("NewAccount() error = %v", err)
		}
		err = a.Get(tests[0].args.id)
		if err != nil {
			t.Errorf("Get() error : %v ", err)
		}
		err = a.Delete()
		if err != nil {
			t.Errorf("Delete() error : %v ", err)
		}
	})

}

func Test_Deposit(t *testing.T) {
	a, err := NewAccount()
	if err != nil {
		t.Fatalf("NewAccount() error : %v ", err)
	}
	const accName = "testacc_76ck76wecoan0vl"

	t.Run("register new account", func(t *testing.T) {
		err = a.Register(accName)
		if err != nil {
			t.Errorf("Register() error : %v ", err)
		}
	})

	var tid int64
	t.Run("deposit account", func(t *testing.T) {
		if err = a.Find(accName); err != nil {
			t.Fatal(err)
		}

		tid, err = a.Deposit(1)
		if err != nil {
			t.Errorf("Deposit() error : %v ", err)
		}
	})

	t.Run("check payment exists", func(t *testing.T) {
		p, err := NewPayment()
		if err != nil {
			t.Fatal(err)
		}
		if err = p.Get(ID(tid)); err != nil {
			t.Errorf("Payment not found: %v", err)
		}
	})

	// delete account
	t.Run("delete account", func(t *testing.T) {
		err = a.Delete()
		if err != nil {
			t.Errorf("Delete() error : %v ", err)
		}
	})

}

func Test_Transfer(t *testing.T) {
	a1, err := NewAccount()
	if err != nil {
		t.Fatalf("NewAccount() error : %v ", err)
	}
	a2, err := NewAccount()
	if err != nil {
		t.Fatalf("NewAccount() error : %v ", err)
	}
	const accName1 = "testacc1_76ck76wecoan0vl"
	const accName2 = "testacc2_76ck76wecoan0vl"

	t.Run("register new account", func(t *testing.T) {
		err = a1.Register(accName1)
		if err != nil {
			t.Fatalf("Register() error : %v ", err)
		}
		err = a2.Register(accName2)
		if err != nil {
			t.Fatalf("Register() error : %v ", err)
		}
	})

	t.Run("deposit account", func(t *testing.T) {
		_, err = a1.Deposit(10)
		if err != nil {
			t.Errorf("Deposit() error : %v ", err)
		}
	})

	// transfer
	var tid int64
	t.Run("transfer", func(t *testing.T) {
		tid, err = a1.Transfer(accName2, 5)
		if err != nil {
			t.Errorf("Transfer() error : %v ", err)
		}
	})

	t.Run("check payment exists", func(t *testing.T) {
		p, err := NewPayment()
		if err != nil {
			t.Fatal(err)
		}
		if err = p.Get(ID(tid)); err != nil {
			t.Errorf("Payment not found: %v", err)
		}
	})

	// check balance
	t.Run("check new balances", func(t *testing.T) {
		_ = a2.Get(a2.ID) // reload from db with new values
		if a1.Balance != 5 || a2.Balance != 5 {
			t.Errorf("New balance error: %v, %v", a1.Balance, a2.Balance)
		}
	})

	// delete accounts
	t.Run("delete account", func(t *testing.T) {
		err = a1.Delete()
		if err != nil {
			t.Errorf("Delete() error : %v ", err)
		}
		err = a2.Delete()
		if err != nil {
			t.Errorf("Delete() error : %v ", err)
		}
	})

}
