//
// Run this test when application is runing!

package main

import (
	"coinswallet/internal/domain/wallet/entity"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"math/rand"
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
			assert.Equal(t, 200, resp.StatusCode())

			var r map[string]string
			if err = json.Unmarshal(resp.Body(), &r); err != nil {
				t.Fatal(err)
			}
			if e, ok := r["error"]; ok {
				t.Error(e)
			}
		}
	})

	t.Run("deposit 4 accounts", func(t *testing.T) {
		for _, n := range wallets {
			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody([]byte(fmt.Sprintf(`{"name":"%s","amount":%f}`, n, rand.Float64()*10+10))).
				Patch("http://localhost:8081/account/deposit/")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, 200, resp.StatusCode())

			var r map[string]interface{}
			if err = json.Unmarshal(resp.Body(), &r); err != nil {
				t.Fatal(err, resp)
			}
			if e, ok := r["error"]; ok {
				t.Error(e)
			}
		}
	})

	t.Run("deposit with errors", func(t *testing.T) {
		w := []struct {
			n    string
			a    float64
			code int
		}{
			{
				n:    "wrongAccountName",
				a:    1,
				code: 404,
			},
			{
				n:    wallets[0],
				a:    -1,
				code: 400,
			},
		}
		for _, n := range w {
			send := func(name string, amount float64) (*resty.Response, error) {
				return client.R().
					SetHeader("Content-Type", "application/json").
					SetBody([]byte(fmt.Sprintf(`{"name":"%s","amount":%f}`, name, amount))).
					Patch("http://localhost:8081/account/deposit/")
			}
			resp, err := send(n.n, n.a)
			assert.Equal(t, n.code, resp.StatusCode())
			if err != nil {
				t.Fatal(err)
			}
			var r map[string]interface{}
			if err = json.Unmarshal(resp.Body(), &r); err != nil {
				t.Fatal(err, resp)
			}
			if e, ok := r["error"]; !ok {
				t.Error(e)
			}
		}
	})

	t.Run("transfer", func(t *testing.T) {
		w := []struct {
			nFrom string
			nTo   string
			a     float64
		}{
			{
				nFrom: wallets[0],
				nTo:   wallets[1],
				a:     2,
			},
			{
				nFrom: wallets[2],
				nTo:   wallets[1],
				a:     2,
			},
			{
				nFrom: wallets[2],
				nTo:   wallets[3],
				a:     2,
			},
			{
				nFrom: wallets[3],
				nTo:   wallets[1],
				a:     2,
			},
			{
				nFrom: wallets[1],
				nTo:   wallets[2],
				a:     2,
			},
		}
		for _, n := range w {
			send := func(name, nameTo string, amount float64) (*resty.Response, error) {
				return client.R().
					SetHeader("Content-Type", "application/json").
					SetBody([]byte(fmt.Sprintf(`{"from":"%s","to":"%s", "amount":%f}`, name, nameTo, amount))).
					Patch("http://localhost:8081/account/transfer/")
			}
			resp, err := send(n.nFrom, n.nTo, n.a)
			assert.Equal(t, 200, resp.StatusCode())
			if err != nil {
				t.Fatal(err)
			}
			var r map[string]interface{}
			if err = json.Unmarshal(resp.Body(), &r); err != nil {
				t.Fatal(err, resp)
			}
			if e, ok := r["error"]; ok {
				t.Error(e)
			}
		}
	})

	t.Run("transfer with errors", func(t *testing.T) {
		w := []struct {
			nFrom string
			nTo   string
			a     float64
			code  int
		}{
			{
				nFrom: wallets[0],
				nTo:   wallets[1],
				a:     -1,
				code:  400,
			},
			{
				nFrom: wallets[1],
				nTo:   wallets[1],
				a:     2,
				code:  400,
			},
			{
				nFrom: wallets[2],
				nTo:   wallets[3],
				a:     0,
				code:  400,
			},
			{
				nFrom: wallets[3],
				nTo:   wallets[1],
				a:     2000,
				code:  400,
			},
			{
				nFrom: "wrongAccountName",
				nTo:   wallets[2],
				a:     2,
				code:  404,
			}, {
				nFrom: wallets[2],
				nTo:   "wrongAccountName",
				a:     2,
				code:  404,
			},
		}
		for _, n := range w {
			send := func(name, nameTo string, amount float64) (*resty.Response, error) {
				return client.R().
					SetHeader("Content-Type", "application/json").
					SetBody([]byte(fmt.Sprintf(`{"from":"%s","to":"%s", "amount":%f}`, name, nameTo, amount))).
					Patch("http://localhost:8081/account/transfer/")
			}
			resp, err := send(n.nFrom, n.nTo, n.a)
			assert.Equal(t, n.code, resp.StatusCode())

			if err != nil {
				t.Fatal(err)
			}
			var r map[string]interface{}
			if err = json.Unmarshal(resp.Body(), &r); err != nil {
				t.Fatal(err, resp)
			}
			if e, ok := r["error"]; !ok {
				t.Error(e)
			}
		}
	})

	t.Run("accounts list", func(t *testing.T) {
		resp, err := client.R().Get("http://localhost:8081/accounts/0/3/")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 200, resp.StatusCode())
		var out map[string]interface{}
		if err = json.Unmarshal(resp.Body(), &out); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, len(out["list"].([]interface{})), 3)
		// o,_ := json.MarshalIndent(out,"","\t")
		// fmt.Println(string(o))
	})

	t.Run("all payment list", func(t *testing.T) {
		resp, err := client.R().Get("http://localhost:8081/payments/0/3/")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 200, resp.StatusCode())
		var out map[string]interface{}
		if err = json.Unmarshal(resp.Body(), &out); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, len(out["list"].([]interface{})), 3)
		// o,_ := json.MarshalIndent(out,"","\t")
		// fmt.Println(string(o))
	})

	t.Run("account payment list", func(t *testing.T) {
		resp, err := client.R().Get(fmt.Sprintf("http://localhost:8081/payments/%s/0/3/", wallets[1]))
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 200, resp.StatusCode())
		var out map[string]interface{}
		if err = json.Unmarshal(resp.Body(), &out); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, len(out["list"].([]interface{})), 3)
		// o,_ := json.MarshalIndent(out,"","\t")
		// fmt.Println(string(o))
	})

	t.Run("delete 4 accounts", func(t *testing.T) {
		for _, n := range wallets {
			a, err := entity.NewAccount()
			if err != nil {
				t.Fatal(err)
			}
			if err = a.Find(entity.AccountName(n)); err != nil {
				t.Fatal(err)
			}
			if err = a.Delete(); err != nil {
				t.Error(err)
			}
		}
	})
}
