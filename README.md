# Form3 Take Home Exercise
Martin Nikolov <martoup@gmail.com>

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


## Instructions
The goal of this exercise is to write a client library in Go to access our fake account API, which is provided as a Docker
container in the file `docker-compose.yaml` of this repository. Please refer to the
[Form3 documentation](http://api-docs.form3.tech/api.html#organisation-accounts) for information on how to interact with the API.

If you encounter any problems running the fake account API we would encourage you to do some debugging first,
before reaching out for help.

### The solution is expected to
- Be written in Go
- Contain documentation of your technical decisions
- Implement the `Create`, `Fetch`, `List` and `Delete` operations on the `accounts` resource. Note that filtering of the List operation is not required, but you should support paging
- Be well tested to the level you would expect in a commercial environment. Make sure your tests are easy to read.

#### Docker-compose
 - Add your solution to the provided docker-compose file
 - We should be able to run `docker-compose up` and see your tests run against the provided account API service 

### Please don't
- Use a code generator to write the client library
- Use (copy or otherwise) code from any third party without attribution to complete the exercise, as this will result in the test being rejected
- Use a library for your client (e.g: go-resty). Only test libraries are allowed.
- Implement an authentication scheme
- Implement support for the fields `data.attributes.private_identification`, `data.attributes.organisation_identification`
  and `data.relationships`, as they are omitted in the provided fake account API implementation
  
## How to submit your exercise
- Include your name in the README. If you are new to Go, please also mention this in the README so that we can consider this when reviewing your exercise
- Create a private [GitHub](https://help.github.com/en/articles/create-a-repo) repository, copy the `docker-compose` from this repository
- [Invite](https://help.github.com/en/articles/inviting-collaborators-to-a-personal-repository) @form3tech-interviewer-1 to your private repo
- Let us know you've completed the exercise using the link provided at the bottom of the email from our recruitment team

## License
Copyright 2019-2020 Form3 Financial Cloud

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
