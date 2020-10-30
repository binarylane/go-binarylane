package binarylane

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAccountGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/account", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		response := `
		{ "account": {
			"server_limit": 25,
			"floating_ip_limit": 25,
			"volume_limit": 22,
			"email": "support@binarylane.com.au",
			"uuid": "b6fr89dbf6d9156cace5f3c78dc9851d957381ef",
			"email_verified": true
			}
		}`

		fmt.Fprint(w, response)
	})

	acct, _, err := client.Account.Get(ctx)
	if err != nil {
		t.Errorf("Account.Get returned error: %v", err)
	}

	expected := &Account{ServerLimit: 25, FloatingIPLimit: 25, Email: "support@binarylane.com.au",
		UUID: "b6fr89dbf6d9156cace5f3c78dc9851d957381ef", EmailVerified: true, VolumeLimit: 22}
	if !reflect.DeepEqual(acct, expected) {
		t.Errorf("Account.Get returned %+v, expected %+v", acct, expected)
	}
}

func TestAccountString(t *testing.T) {
	acct := &Account{
		ServerLimit:     25,
		FloatingIPLimit: 25,
		Email:           "support@binarylane.com.au",
		UUID:            "b6fr89dbf6d9156cace5f3c78dc9851d957381ef",
		EmailVerified:   true,
		Status:          "active",
		StatusMessage:   "message",
		VolumeLimit:     22,
	}

	stringified := acct.String()
	expected := `binarylane.Account{ServerLimit:25, FloatingIPLimit:25, VolumeLimit:22, Email:"support@binarylane.com.au", UUID:"b6fr89dbf6d9156cace5f3c78dc9851d957381ef", EmailVerified:true, Status:"active", StatusMessage:"message"}`
	if expected != stringified {
		t.Errorf("Account.String returned %+v, expected %+v", stringified, expected)
	}

}
