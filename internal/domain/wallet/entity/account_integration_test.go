// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// yuriyiovkov@gmail.com; telegram: @yuriyiovkov

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
				"mywallet",
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
				"mywallet",
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
		t.Errorf("NewAccount() error : %v ", err)
		t.Fail()
	}
	const accName = "testacc"

	t.Log("register new account")
	err = a.Register(accName)
	if err != nil {
		t.Errorf("Register() error : %v ", err)
	}

	t.Log("deposit account")
	if err = a.Find(accName); err != nil {
		t.Fatal(err)
	}

	tid, err := a.Deposit(1)
	if err != nil {
		t.Errorf("Deposit() error : %v ", err)
	}
	t.Log("payment id: ", tid)

	_ = a.Delete()

	// TODO
	// check that payment with id exists
	_ = tid

}
