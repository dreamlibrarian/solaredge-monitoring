package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
)

/*

arright, invocation model abstractions suck.

AWS credentials will come from the lambda policy.

API key needs to come from an AWS vault.
Environment should specify the vault and key name.

Prrrrobably want to have some kind of checkpointing mechanic.

The rest of it should be in the message:
verb, arguments like the cmd entrypoint.





*/

// var client = lambda.New(session.New())

func handleRequest(ctx context.Context, event events.SQSEvent) (string, error) {
	/*
		lctx := lambdacontext.FromContext(ctx)

		for _, record := range event.Records {

		}
	*/

	return "", nil
}

func actionForRecord(message events.SQSMessage) (string, error) {

	return "", nil
}

func main() {
	runtime.Start(handleRequest)
}
