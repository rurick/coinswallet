// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// yuriyiovkov@gmail.com; telegram: @yuriyiovkov

package entity

import (
	"testing"
)

func Test_List(t *testing.T) {
	var (
		err error
		lst []Payment
	)

	t.Run("get all payments", func(t *testing.T) {
		lst, err = PaymentsList(nil, 0, -1)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("get account payments", func(t *testing.T) {
		ac, _ := NewAccount()
		_ = ac.Register("random_76ck76wecoan0vl12")
		_, _ = ac.Deposit(1)
		lst, err = PaymentsList(ac, 0, -1)
		if len(lst) != 1 {
			t.Errorf("wait 1 row got: %d", len(lst))
		}
		_ = ac.Delete()
		if err != nil {
			t.Fatal(err)
		}
	})
}
