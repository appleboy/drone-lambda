package main

import (
	"errors"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/sirupsen/logrus"
)

type (
	// Config for the plugin.
	Config struct {
		Region          string
		AccessKey       string
		SecretKey       string
		Profile         string
		FunctionName    string
		S3Bucket        string
		S3Key           string
		S3ObjectVersion string
		ZipFile         string
	}

	// Plugin values.
	Plugin struct {
		Config Config
	}
)

// Exec executes the plugin.
func (p Plugin) Exec() error {
	if p.Config.FunctionName == "" {
		return errors.New("missing lambda function name")
	}

	// Create Lambda service client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	config := &aws.Config{
		Region: aws.String(p.Config.Region),
	}

	if p.Config.Profile != "" {
		config.Credentials = credentials.NewSharedCredentials("", p.Config.Profile)
	}

	input := &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(p.Config.FunctionName),
		Publish:      aws.Bool(true),
	}

	if p.Config.S3Bucket != "" && p.Config.S3Key != "" {
		input.S3Bucket = aws.String(p.Config.S3Bucket)
		input.S3Key = aws.String(p.Config.S3Key)

		if p.Config.S3ObjectVersion != "" {
			input.S3ObjectVersion = aws.String(p.Config.S3ObjectVersion)
		}
	}

	if p.Config.ZipFile != "" {
		contents, err := ioutil.ReadFile(p.Config.ZipFile)

		if err != nil {
			logrus.Println("Could not read " + p.Config.ZipFile)
			return err
		}

		input.ZipFile = contents
	}

	svc := lambda.New(sess, config)

	result, err := svc.UpdateFunctionCode(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case lambda.ErrCodeServiceException:
				logrus.Println(lambda.ErrCodeServiceException, aerr.Error())
			case lambda.ErrCodeResourceNotFoundException:
				logrus.Println(lambda.ErrCodeResourceNotFoundException, aerr.Error())
			case lambda.ErrCodeInvalidParameterValueException:
				logrus.Println(lambda.ErrCodeInvalidParameterValueException, aerr.Error())
			case lambda.ErrCodeTooManyRequestsException:
				logrus.Println(lambda.ErrCodeTooManyRequestsException, aerr.Error())
			case lambda.ErrCodeCodeStorageExceededException:
				logrus.Println(lambda.ErrCodeCodeStorageExceededException, aerr.Error())
			default:
				logrus.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			logrus.Println(err.Error())
		}
		return err
	}

	logrus.Println(result)

	return nil
}
