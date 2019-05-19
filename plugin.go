package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/mholt/archiver"
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
		DryRun          bool
		ZipFile         string
		Source          []string
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

	if p.Config.S3Bucket == "" &&
		p.Config.S3Key == "" &&
		len(p.Config.Source) == 0 &&
		p.Config.ZipFile == "" {
		return errors.New("missing zip source")
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

	if p.Config.AccessKey != "" && p.Config.SecretKey != "" {
		config.Credentials = credentials.NewStaticCredentials(p.Config.AccessKey, p.Config.SecretKey, "")
	}

	input := &lambda.UpdateFunctionCodeInput{
		DryRun:       aws.Bool(p.Config.DryRun),
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

	if len(p.Config.Source) != 0 {
		files := globList(trimPath(p.Config.Source))
		path := os.TempDir() + "/output.zip"
		zip := archiver.NewZip()
		if len(files) != 0 {
			if err := zip.Archive(files, path); err != nil {
				log.Fatalf("can't create zip file: %s", err.Error())
			}

			p.Config.ZipFile = path
		}
	}

	if p.Config.ZipFile != "" {
		contents, err := ioutil.ReadFile(p.Config.ZipFile)

		if err != nil {
			log.Println("Could not read " + p.Config.ZipFile)
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
				log.Println(lambda.ErrCodeServiceException, aerr.Error())
			case lambda.ErrCodeResourceNotFoundException:
				log.Println(lambda.ErrCodeResourceNotFoundException, aerr.Error())
			case lambda.ErrCodeInvalidParameterValueException:
				log.Println(lambda.ErrCodeInvalidParameterValueException, aerr.Error())
			case lambda.ErrCodeTooManyRequestsException:
				log.Println(lambda.ErrCodeTooManyRequestsException, aerr.Error())
			case lambda.ErrCodeCodeStorageExceededException:
				log.Println(lambda.ErrCodeCodeStorageExceededException, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return err
	}

	log.Println(result)

	return nil
}

func trimPath(keys []string) []string {
	var newKeys []string

	for _, value := range keys {
		value = strings.Trim(value, " ")
		if len(value) == 0 {
			continue
		}

		newKeys = append(newKeys, value)
	}

	return newKeys
}

func globList(paths []string) []string {
	var newPaths []string

	for _, pattern := range paths {
		pattern = strings.Trim(pattern, " ")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			log.Printf("Glob error for %q: %s\n", pattern, err)
			continue
		}

		newPaths = append(newPaths, matches...)
	}

	return newPaths
}
