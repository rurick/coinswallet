// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// yuriyiovkov@gmail.com; telegram: @yuriyiovkov

package wallet

import (
	"testing"
)

func Test_Validate(t *testing.T) {
	type args struct {
		name AccountName
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"valid mixed-case account name",
			args{
				"QWEasdf45",
			},
			false,
		},
		{
			"valid lowercase account name",
			args{
				"asdf45",
			},
			false,
		},
		{
			"invalid account name",
			args{
				"asdf 45",
			},
			true,
		},
		{
			"empty account name",
			args{
				"",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wa := Account{}
			err := wa.Validate(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v ", err, tt.wantErr)
			}
		})
	}
}
