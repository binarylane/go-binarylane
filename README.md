# go-binarylane

[![Build Status](https://travis-ci.org/binarylane/go-binarylane.svg)](https://travis-ci.org/binarylane/go-binarylane)
[![GoDoc](https://godoc.org/github.com/binarylane/go-binarylane?status.svg)](https://godoc.org/github.com/binarylane/go-binarylane)

go-binarylane is a Go client library for accessing the BinaryLane API.

You can view the client API docs here: [http://godoc.org/github.com/binarylane/go-binarylane](http://godoc.org/github.com/binarylane/go-binarylane)

You can view BinaryLane API docs here: [https://api.binarylane.com.au/reference/](https://api.binarylane.com.au/reference/)

## Install
```sh
go get github.com/binarylane/go-binarylane@vX.Y.Z
```

where X.Y.Z is the [version](https://github.com/binarylane/go-binarylane/releases) you need.

or
```sh
go get github.com/binarylane/go-binarylane
```
for non Go modules usage or latest version.

## Usage

```go
import "github.com/binarylane/go-binarylane"
```

Create a new BinaryLane client, then use the exposed services to
access different parts of the BinaryLane API.

### Authentication

Currently, Access Token is the only method of
authenticating with the API. You can manage your tokens
at the BinaryLane Control Panel [Developer API](https://home.binarylane.com.au/api-info).

You can then use your token to create a new client:

```go
package main

import (
    "github.com/binarylane/go-binarylane"
)

func main() {
    client := binarylane.NewFromToken("my-binarylane-api-token")
}
```

If you need to provide a `context.Context` to your new client, you should use [`binarylane.NewClient`](https://godoc.org/github.com/binarylane/go-binarylane#NewClient) to manually construct a client instead.

## Examples


To create a new Server:

```go
serverName := "my-first-server"

createRequest := &binarylane.DropletCreateRequest{
    Name:   serverName,
    Region: "syd",
    Size:   "std-1vcpu",
    Image: binarylane.ServerCreateImage{
        Slug: "ubuntu-20-04-lts",
    },
}

ctx := context.TODO()

newServer, _, err := client.Servers.Create(ctx, createRequest)

if err != nil {
    fmt.Printf("Something bad happened: %s\n\n", err)
    return err
}
```

### Pagination

If a list of items is paginated by the API, you must request pages individually. For example, to fetch all Servers:

```go
func ServerList(ctx context.Context, client *binarylane.Client) ([]binarylane.Server, error) {
    // create a list to hold our servers
    list := []binarylane.Server{}

    // create options. initially, these will be blank
    opt := &binarylane.ListOptions{}
    for {
        droplets, resp, err := client.Servers.List(ctx, opt)
        if err != nil {
            return nil, err
        }

        // append the current page's servers to our list
        list = append(list, servers...)

        // if we are at the last page, break out the for loop
        if resp.Links == nil || resp.Links.IsLastPage() {
            break
        }

        page, err := resp.Links.CurrentPage()
        if err != nil {
            return nil, err
        }

        // set the page we want for the next request
        opt.Page = page + 1
    }

    return list, nil
}
```

## Versioning

Each version of the client is tagged and the version is updated accordingly.

To see the list of past versions, run `git tag`.


## Documentation

For a comprehensive list of examples, check out the [API documentation](https://api.binarylane.com.au/reference/).

For details on all the functionality in this library, see the [GoDoc](http://godoc.org/github.com/binarylane/go-binarylane) documentation.


## Contributing

We love pull requests! Please see the [contribution guidelines](CONTRIBUTING.md).
