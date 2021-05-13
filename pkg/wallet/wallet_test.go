package wallet

import (
	"context"
	"testing"
)

func Test_Validate(t *testing.T) {
	type args struct {
		ctx  context.Context
		name AccountName
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"valid mixedcase account name",
			args{
				context.Background(),
				"QWEasdf45",
			},
			false,
		},
		{
			"valid lowercase account name",
			args{
				context.Background(),
				"asdf45",
			},
			false,
		},
		{
			"invalid account name",
			args{
				context.Background(),
				"asdf 45",
			},
			true,
		},
		{
			"empty account name",
			args{
				context.Background(),
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
