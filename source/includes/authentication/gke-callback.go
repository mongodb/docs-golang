package main

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// start-callback
	gkeCallback := func(_ context.Context,
		_ *options.OIDCArgs) (*options.OIDCCredential, error) {
		accessToken, err := os.ReadFile(
			"/var/run/secrets/kubernetes.io/serviceaccount/token")
		if err != nil {
			return nil, err
		}
		return &options.OIDCCredential{
			AccessToken: string(accessToken),
		}, nil
	}
	// end-callback

	// start-credential-callback
	uri := "mongodb://<hostname>:<port>"
	props := map[string]string{
		"ENVIRONMENT":    "gcp",
		"TOKEN_RESOURCE": "<audience>",
	}
	opts := options.Client().ApplyURI(uri)
	opts.SetAuth(
		options.Credential{
			AuthMechanism:           "MONGODB-OIDC",
			AuthMechanismProperties: props,
			OIDCMachineCallback:     gkeCallback,
		},
	)
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	// end-credential-callback
}