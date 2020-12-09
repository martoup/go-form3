package form3

import (
	"context"
	"fmt"
	"net/http"
)

const (
	// accountsPath URL path to accounts resources.
	accountsPath = "organisation/accounts"
)

// Account represents a single bank account that is registered with Form3.
type Account struct {
	Data *AccountData `json:"data"`
}

// AccountList represents a list of bank accounts that is registered with Form3.
type AccountList struct {
	Data []*AccountData `json:"data"`
}

// AccountData represents the main attributes for a given Form3 account.
type AccountData struct {
	Type           string             `json:"type"`
	ID             string             `json:"id"`
	OrganisationID string             `json:"organisation_id"`
	Version        int                `json:"version"`
	Attributes     *AccountAttributes `json:"attributes"`
	CreatedOn      string             `json:"created_on,omitempty"`
	ModifiedOn     string             `json:"modified_on,omitempty"`
}

// AccountAttributes represents the available attribute fields.
// The availability of each field depends on the API call and scheme.
type AccountAttributes struct {
	Country                 string   `json:"country"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	CustomerID              string   `json:"customer_id,omitempty"`
	Name                    []string `json:"name"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	AccountClassification   string   `json:"account_classification,omitempty"`
	JointAccount            bool     `json:"joint_account,omitempty"`
	AccountMatchingOptOut   bool     `json:"account_matching_opt_out,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Switched                bool     `json:"switched,omitempty"`
	Status                  string   `json:"status"`
}

// AccountsService handles the communication with the account related
// methods of the Form3 API.
//
// Form3 API docs: https://api-docs.form3.tech/api.html?http#organisation-accounts
type AccountsService service

// Create registers an existing bank account with Form3 or create a new one.
// The country attribute must be specified as a minimum.
// Depending on the country, other attributes such as bank_id and bic are mandatory.
func (s *AccountsService) Create(ctx context.Context, account *Account) (*Account, *http.Response, error) {
	request, err := s.client.NewRequest(http.MethodPost, accountsPath, account)
	if err != nil {
		return nil, nil, err
	}

	acc := new(Account)
	resp, err := s.client.Do(ctx, request, acc)
	if err != nil {
		return nil, resp, err
	}

	return acc, resp, nil
}

// Fetch gets a single account using the account ID.
func (s *AccountsService) Fetch(ctx context.Context, accountID string) (*Account, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", accountsPath, accountID)
	request, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	acc := new(Account)
	resp, err := s.client.Do(ctx, request, acc)
	if err != nil {
		return nil, resp, err
	}

	return acc, resp, nil
}

// List lists all accounts. Supports pagination.
func (s *AccountsService) List(ctx context.Context, pageNumber int, pageSize int) (*AccountList, *http.Response, error) {
	path := fmt.Sprintf("%s?page[number]=%d&page[size]=%d", accountsPath, pageNumber, pageSize)
	request, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	acc := new(AccountList)
	resp, err := s.client.Do(ctx, request, &acc)
	if err != nil {
		return nil, resp, err
	}

	return acc, resp, nil
}

// Delete deletes an account by ID and given version
func (s *AccountsService) Delete(ctx context.Context, accountID string, version int) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s?version=%d", accountsPath, accountID, version)
	request, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, request, nil)
}
