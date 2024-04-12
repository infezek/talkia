package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

func login() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Login to the chat",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			opt := option.WithCredentialsFile("./firebase.json")
			app, err := firebase.NewApp(ctx, nil, opt)
			if err != nil {
				log.Fatalf("error initializing app: %v\n", err)
			}
			client, err := app.Auth(ctx)
			if err != nil {
				log.Fatalf("error getting Auth client: %v\n", err)
			}
			userRecord, err := client.GetUserByProviderUID(ctx, "email", "ezequiel.jn98@outlook.com")
			if err != nil {
				log.Fatalf("error getting user record: %v\n", err)
			}
			fmt.Printf("User UID: %s\n", userRecord)
		},
	}

}
