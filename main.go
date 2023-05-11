package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

func main() {
	http.HandleFunc("/getsecret", SecretsVault)
	http.ListenAndServe(":8080", nil)
}

func initVault() (*vault.Client, context.Context) {
	ctx := context.Background()
	client, err := vault.New(
		vault.WithEnvironment(),
	); if err != nil {
		log.Fatal(err)
	}

	defaultRequest := schema.NewAzureLoginRequestWithDefaults()
	_, err = client.Auth.AzureLogin(
		ctx,
		*defaultRequest,
	); if err != nil {
		log.Fatal(err)
	}

	return client, ctx
}

func SecretsVault(w http.ResponseWriter, r *http.Request) {
	client, ctx := initVault()
	s, err := client.Secrets.KvV2Read(ctx, "my-secrets/poc-interbank"); if err != nil {
		log.Fatalf("Error when reading the secret: %s", err)
	}
	fmt.Fprintf(w, "Your Secret is:%s", s.Data.Data)
}