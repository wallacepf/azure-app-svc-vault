package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

func main() {
	http.HandleFunc("/getsecret", secretsVault)
	http.ListenAndServe(":8080", nil)
}

func initVault() (*vault.Client, context.Context) {
	ctx := context.Background()
	// clientID := azidentity.ClientID("18aad86a-2514-482e-a7d4-b56d5890eece")
	// opts := azidentity.ManagedIdentityCredentialOptions{ID: clientID}
	cred, err := azidentity.NewDefaultAzureCredential(nil); if err != nil {
	// cred, err := azidentity.NewManagedIdentityCredential(&opts); if err != nil {
		log.Printf("Error on getting Azure Token: %s", err)
	}
	var scopes []string
	scopes = append(scopes, "https://management.azure.com/")

	tokReqOptions := policy.TokenRequestOptions{}
	tokReqOptions.Scopes = scopes

	token, err := cred.GetToken(ctx, tokReqOptions); if err != nil {
		log.Fatalf("Unable to retrieve the token: %s", err)
	}

	client, err := vault.New(
		vault.WithEnvironment(),
		// vault.WithAddress("http://127.0.0.1:8200"),
	); if err != nil {
		log.Fatalf("Cannot connect to the Vault Instance: %s", err)
	}

	defaultRequest := schema.AzureLoginRequest{
		Jwt: token.Token,
		ResourceGroupName: "MyDemoAppRG",
		Role: "myapp",
		SubscriptionId: "7f7602dd-85a6-4140-8501-61f2ee9f65a9",
		ResourceId: "/subscriptions/7f7602dd-85a6-4140-8501-61f2ee9f65a9/resourcegroups/MyDemoAppRG/providers/Microsoft.Web/sites/myapp-demo-pov",
	}

	fmt.Print(token.Token)
	_, err = client.Auth.AzureLogin(
		ctx,
		defaultRequest,
	); if err != nil {
		log.Fatalf("Error loggin on Vault with Azure Creds: %s", err)
	}

	return client, ctx
}

func secretsVault(w http.ResponseWriter, r *http.Request) {
	client, ctx := initVault()
	s, err := client.Secrets.KvV2Read(ctx, "my-secrets/poc"); if err != nil {
		log.Fatalf("Error when reading the secret: %s", err)
	}
	fmt.Fprintf(w, "Your Secret is:%s", s.Data.Data)
}