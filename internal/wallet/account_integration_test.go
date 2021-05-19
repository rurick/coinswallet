// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// yuriyiovkov@gmail.com; telegram: @yuriyiovkov

package wallet

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
			"with to short name",
			args{
				"asf",
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v ", err, tt.wantErr)
				return
			}
			err = a.Register(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v ", err, tt.wantErr)
				return
			}
			tt.args.id = a.ID
		})
	}

	// delete created accounts
	t.Run("delete existing accounts", func(t *testing.T) {
		a, _ := New()
		err := a.Get(tests[0].args.id)
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
	a, err := New()
	if err != nil {
		t.Errorf("New() error : %v ", err)
	}

	t.Run("register account", func(t *testing.T) {
		err := a.Register("testacc")
		if err != nil {
			t.Errorf("Register() error : %v ", err)
		}
	})

	t.Run("deposit account", func(t *testing.T) {
		tid, err := a.Deposit(1)
		if err != nil {
			t.Errorf("Deposit() error : %v ", err)
		}

		// TODO
		// check that payment with id exists
		_ = tid
	})

}
