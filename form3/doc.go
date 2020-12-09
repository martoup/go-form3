/*
Package form3 provides a client for using the Form3 API.

Usage:

	import "github.com/martoup/go-form3/form3"

Construct a new Form3 client, then use the account service on the client to
access the Form3 API. For example:

	client, err := form3.NewClientFromEnvironment()

	if err != nil {
		log.Fatalf("Failed to create a client %v", err)
	}

	// create an account
	create, _, err := client.Accounts.Create(ctx, account)

    // get a single account by ID
	fetch, _, err := client.Accounts.Fetch(ctx, accountId)

    // list all accounts
	list, _, err := client.Accounts.List(ctx, 0, 0)

    // delete an account
	_, err = client.Accounts.Delete(ctx, accountId, accountVersion)
*/
package form3
