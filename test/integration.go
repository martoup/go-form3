package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/martoup/go-form3/form3"
	"log"
	"reflect"
	"strings"
)

const accJson = `
{
  "data": {
    "type": "accounts",
    "id": "1227e265-9605-4b4b-a0e5-3003ea9cc4dc",
    "organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
    "attributes": {
      "country": "GB",
      "base_currency": "GBP",
      "bank_id": "400300",
      "bank_id_code": "GBDSC",
      "bic": "NWBKGB22",
      "account_classification": "Personal",
      "joint_account": false,
      "account_matching_opt_out": false,
      "secondary_identification": "A1B2C3D4"
    }
  }
}`

func main() {
	client, err := form3.NewClientFromEnvironment()

	if err != nil {
		log.Fatalf("Failed to create a client %v", err)
	}

	account := &form3.Account{}
	err = json.NewDecoder(strings.NewReader(accJson)).Decode(account)

	if err != nil {
		log.Fatalf("Failed to decode test JSON %v", err)
	}

	fmt.Printf("==== Step 1/5 Create account: %v\n", accJson)
	create, _, err := client.Accounts.Create(context.Background(), account)

	if err != nil {
		log.Fatalf("Failed to create account %v", err)
	}

	fmt.Printf("Response: %+v\n\n", printJSON(create))
	checkJSON(create)

	fmt.Printf("==== Step 2/5 Get single account with ID %s\n", account.Data.ID)
	fetch, _, err := client.Accounts.Fetch(context.Background(), account.Data.ID)

	if err != nil {
		log.Fatalf("Failed to fetch account %v\n", err)
	}
	fmt.Printf("Response: %+v\n\n", printJSON(fetch))
	checkJSON(fetch)

	fmt.Print("==== Step 3/5 Get account list:\n")
	list, _, err := client.Accounts.List(context.Background(), 0, 0)

	if err != nil {
		log.Fatalf("Failed to get account list %v", err)
	}

	fmt.Printf("Response: %+v\n\n", printJSON(list))

	fmt.Printf("==== Step 4/5 Delete account with ID: %s and Version: %d \n", account.Data.ID, account.Data.Version)
	_, err = client.Accounts.Delete(context.Background(), account.Data.ID, account.Data.Version)

	if err != nil {
		log.Fatalf("Failed to delete account %v", err)
	}

	fmt.Print("==== Step 5/5 Get single (now deleted) account:\n")
	fetch, _, err = client.Accounts.Fetch(context.Background(), account.Data.ID)

	if err == nil {
		log.Fatal(err)
	}

	fmt.Printf("Response should be 404 (we deleted the account), actual response is: %+v\n", err)
	fmt.Print("===== Success! =====")
}

func printJSON(body interface{}) string {
	marshal, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Failed to marshal body to print %v", err)
	}
	return string(marshal)
}

func checkJSON(act *form3.Account) {
	fmt.Print("Checking Response... \n")
	exp := form3.Account{}
	err := json.Unmarshal([]byte(accJson), &exp)
	if err != nil {
		log.Fatalf("Failed to marshal body %v", err)
	}

	if !reflect.DeepEqual(act.Data.Attributes, exp.Data.Attributes) {
		log.Fatalf("Objects do not match %+v, expected %+v", act, exp)
	}
	fmt.Print("Asserting Objects successful.\n")
}
