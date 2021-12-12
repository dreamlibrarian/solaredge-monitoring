package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/dreamlibrarian/solaredge-monitoring/action"
	"github.com/dreamlibrarian/solaredge-monitoring/api"
	"github.com/rs/zerolog/log"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
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

var client = lambda.New(session.New())

const (
	envBucketNameKey   = "bucketName"
	envBucketPrefixKey = "bucketPrefix"
	envTimeUnitKey     = "timeUnit"
	envSiteIDsKey      = "siteIDs"
	envSecretID        = "apiKey"

	checkpointKey = "solaredge-monitoring-checkpoint"
)

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

func energyAction(ctx context.Context, message events.SQSMessage) error {
	lctx, isLambdaContext := lambdacontext.FromContext(ctx)
	if !isLambdaContext {
		return errors.New("invoked energyAction from non-lambda context, no idea how to proceed")
	}
	config := &action.EnergyConfig{}

	env := lctx.ClientContext.Env

	if timeUnit, ok := env[envTimeUnitKey]; ok {
		config.TimeUnit = timeUnit
	}
	if siteIDs, ok := env[envSiteIDsKey]; ok {
		config.SiteIDs = strings.Split(",", siteIDs)
	} else {
		config.DiscoverSites = true
	}

	checkpoint, err := getCheckpoint(ctx)
	if err != nil {
		return err
	}

	config.StartTime = checkpoint

	apiKey, err := getAPIKey(ctx)
	if err != nil {
		return err
	}

	log.Debug().Interface("actionConfig", config).Msg("Invoking energy endpoint")

	act := action.NewEnergyAction(apiKey)

	result, err := act.Do(config)
	if err != nil {
		return err
	}

	return nil
}

func getAPIKey(ctx context.Context) (string, error) {
	lctx, isLambdaContext := lambdacontext.FromContext(ctx)
	if !isLambdaContext {
		return "", errors.New("invoked energyAction from non-lambda context, no idea how to proceed")
	}
	env := lctx.ClientContext.Env

	secretID, ok := env[envSecretID]
	if !ok {
		return "", errors.New("unable to fetch API key name from environment")
	}

	sess, err := session.NewSession()
	if err != nil {
		return "", err
	}

	sm := secretsmanager.New(sess)
	if err != nil {
		return "", err
	}

	secret, err := sm.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	})
	if err != nil {
		return "", err
	}

	return secret.String(), nil

}

func getCheckpoint(ctx context.Context) (time.Time, error) {
	lctx, isLambdaContext := lambdacontext.FromContext(ctx)
	if !isLambdaContext {
		return time.Time{}, errors.New("invoked energyAction from non-lambda context, no idea how to proceed")
	}
	env := lctx.ClientContext.Env

	bucketName, ok := env[envBucketNameKey]
	if !ok {
		return time.Time{}, errors.New("no bucket name specified in config")
	}
	bucketPrefix, _ := env[envBucketNameKey]

	sess, err := session.NewSession()
	if err != nil {
		return time.Time{}, err
	}

	s3Client := s3.New(sess)

	object, err := s3Client.GetObject(&s3.GetObjectInput{
		Key:    aws.String(fmt.Sprintf("%s/%s", bucketPrefix, checkpointKey)),
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "NotFound" {
				// no checkpoint? Assume everything's fine and start from defaults.
				return time.Time{}, nil
			}
		}
		return time.Time{}, err
	}

	return api.ParseTime(object.String())
}

func setCheckpoint(ctx context.Context, t time.Time) error {
	lctx, isLambdaContext := lambdacontext.FromContext(ctx)
	if !isLambdaContext {
		return errors.New("invoked energyAction from non-lambda context, no idea how to proceed")
	}
	env := lctx.ClientContext.Env

	bucketName, ok := env[envBucketNameKey]
	if !ok {
		return errors.New("no bucket name specified in config")
	}

	bucketPrefix, _ := env[envBucketPrefixKey]

	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Body:   strings.NewReader(api.ToTimestamp(t)),
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", bucketPrefix, checkpointKey)),
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	runtime.Start(handleRequest)
}
