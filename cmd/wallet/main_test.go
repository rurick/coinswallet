package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_API(t *testing.T) {
	wallets := []string{
		"testwalletsalt1",
		"testwalletsalt2",
		"testwalletsalt3",
		"testwalletsalt4",
	}
	client := resty.New()

	t.Run("create 4 accounts", func(t *testing.T) {
		for _, n := range wallets {
			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody([]byte(fmt.Sprintf(`{"name":"%s"}`, n))).
				Post("http://localhost:8081/account/")
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(resp)
		}

	})

	resp, err := client.R().Get("http://localhost:8081/accounts/0/-1/")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 200, resp.StatusCode())
	fmt.Println(resp)
}
