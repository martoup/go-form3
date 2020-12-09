package form3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var expectedAccount = &Account{Data: &AccountData{
	Type:           "accounts",
	ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
	OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
	Version:        0,
	Attributes: &AccountAttributes{
		Country:                 "GB",
		BaseCurrency:            "GBP",
		BankID:                  "400300",
		BankIDCode:              "GBDSC",
		Bic:                     "NWBKGB22",
		Name:                    []string{"Samantha Holder"},
		AlternativeNames:        []string{"Sam Holder"},
		AccountClassification:   "Personal",
		JointAccount:            false,
		AccountMatchingOptOut:   false,
		SecondaryIdentification: "A1B2C3D4",
	},
}}

func TestAccountsService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/"+accountsPath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		response, err := json.Marshal(expectedAccount)
		testBody(t, r, bytes.NewBuffer(response))
		if err != nil {
			t.Errorf("Unexpected error in test data: %v", err)
		}
		fmt.Fprint(w, string(response))
	})

	acct, _, err := client.Accounts.Create(ctx, expectedAccount)
	if err != nil {
		t.Errorf("Accounts.Create returned error: %v", err)
	}

	if !reflect.DeepEqual(acct, expectedAccount) {
		t.Errorf("Accounts.Create returned %+v, expected %+v", acct, expectedAccount)
	}
}

func TestAccountsService_Fetch(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/"+accountsPath+"/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		response, err := json.Marshal(expectedAccount)
		if err != nil {
			t.Errorf("Unexpected error in test data: %v", err)
		}
		fmt.Fprint(w, string(response))
	})

	acct, _, err := client.Accounts.Fetch(ctx, "1")
	if err != nil {
		t.Errorf("Accounts.Fetch returned error: %v", err)
	}

	if !reflect.DeepEqual(acct, expectedAccount) {
		t.Errorf("Accounts.Fetch returned %+v, expected %+v", acct, expectedAccount)
	}
}

func TestAccountsService_List(t *testing.T) {
	setup()
	defer teardown()

	accountListResponse := AccountList{Data: []*AccountData{expectedAccount.Data}}

	mux.HandleFunc("/v1/"+accountsPath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testQueryParam(t, r, "page[number]", "0")
		testQueryParam(t, r, "page[size]", "0")

		response, err := json.Marshal(accountListResponse)
		if err != nil {
			t.Errorf("Unexpected error in test data: %v", err)
		}
		fmt.Fprint(w, string(response))
	})

	acct, _, err := client.Accounts.List(ctx, 0, 0)
	if err != nil {
		t.Errorf("Accounts.List returned error: %v", err)
	}

	if !reflect.DeepEqual(acct.Data[0], accountListResponse.Data[0]) {
		t.Errorf("Accounts.List returned %+v, expected %+v", acct, accountListResponse)
	}
}

func TestAccountsService_ListEmpty(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/"+accountsPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"data\":[]}")
	})

	acct, _, err := client.Accounts.List(ctx, 0, 0)
	if err != nil {
		t.Errorf("Accounts.List returned error: %v", err)
	}

	if len(acct.Data) != 0 {
		t.Errorf("Account Data should be empty.")
	}
}

func TestAccountsService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/"+accountsPath+"/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testQueryParam(t, r, "version", "0")
	})

	_, err := client.Accounts.Delete(ctx, "1", 0)
	if err != nil {
		t.Errorf("Accounts.Delete returned error: %v", err)
	}
}

func TestAccountsService_DeleteNotFound(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/"+accountsPath+"/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	resp, err := client.Accounts.Delete(ctx, "1", 0)

	if err == nil {
		t.Errorf("Accounts.Delete should return an error on 404")
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Response Code not correct. Expected %v got %v", resp.StatusCode, resp.StatusCode)
	}
}
