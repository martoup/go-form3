## About
This project contains go client library for the [Form3 Accounts API](https://api-docs.form3.tech/api.html?http#organisation-accounts).
The source code is located in the `form3` directory. In addition, a Dockerfile is created and a new `integration-test` service is added to the `docker-compose.yml` 

I have split the code into two main files:
 - `form3.go` containing a `Client` - a thin wrapper around standard `http.Client` 
(I haven't used a third-party rest library) that abstracts two methods `NewRequest()` and `Do()`. 
They contain common http logic later used in `AccountService`.
- `accounts.go` containing the structures and the required API operations.

I like this code architecture of separating the code in one abstract `Client` and specific `Services` that use it. 
In this way we avoid duplicating code related to http communication, we have clear separation of 
responsibilities and the API client is easily extensible. 

I also like this style of writing go client libraries because is widely used, recognisable and adopted in some of the 
most-popular go libs like [go-github](https://github.com/google/go-github) and [godo](https://github.com/digitalocean/godo)

I added a ground layer of unit tests, which catch most of the cases. 
In a real-world project even more extensive test suite could be implemented. In addition - in the `/test` folder a sample 
script could be found. It covers all Happy-path API requests and it is also used in the `docker-compose` file.   
  

## Usage
Run a sample demo with:
```
docker-compose up
```

Run tests with:
```
go test ./...
```

Construct a new Form3 client, then use the account service on the client to
access the Form3 API. For example:
```
import "github.com/martoup/go-form3/form3"

// Make sure FORM3_BASE_URL is set.
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
```
_Other examples can be found in `/test/integration.go`_ 
