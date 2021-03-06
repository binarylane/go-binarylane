package util

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/binarylane/go-binarylane"
)

func ExampleWaitForActive() {
	// build client
	pat := "mytoken"
	token := &oauth2.Token{AccessToken: pat}
	t := oauth2.StaticTokenSource(token)

	ctx := context.TODO()
	oauthClient := oauth2.NewClient(ctx, t)
	client := binarylane.NewClient(oauthClient)

	// create your server and retrieve the create action uri
	uri := "https://api.binarylane.com.au/v2/actions/xxxxxxxx"

	// block until until the action is complete
	err := WaitForActive(ctx, client, uri)
	if err != nil {
		panic(err)
	}
}
