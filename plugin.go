package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/gookit/goutil/dump"
	"github.com/mholt/archiver/v3"
)

type (
	// Config for the plugin.
	Config struct {
		Region          string
		AccessKey       string
		SecretKey       string
		Profile         string
		FunctionName    string
		ReversionID     string
		S3Bucket        string
		S3Key           string
		S3ObjectVersion string
		DryRun          bool
		ZipFile         string
		Source          []string
		Debug           bool
		Publish         bool
		MemorySize      int64
		Timeout         int64
		Handler         string
		Role            string
		Runtime         string
		Environment     []string
		ImageURI        string
		Subnets         []string
		SecurityGroups  []string
		Description     string
		Layers          []string
		SessionToken    string
		TracingMode     string
		MaxAttempts     int
	}

	// Commit information.
	Commit struct {
		Sha    string
		Author string
	}

	// Plugin values.
	Plugin struct {
		Config Config
		Commit Commit
	}
)

func getEnvironment(Environment []string) map[string]string {
	output := make(map[string]string)
	for _, e := range Environment {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) != 2 {
			continue
		}
		output[pair[0]] = pair[1]
	}
	return output
}

func (p Plugin) loadEnvironment(envs []string) *lambda.Environment {
	return &lambda.Environment{
		Variables: aws.StringMap(getEnvironment(envs)),
	}
}

// Exec executes the plugin.
func (p Plugin) Exec() error {
	p.dump(p.Config)

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
		config.Credentials = credentials.NewStaticCredentials(p.Config.AccessKey, p.Config.SecretKey, p.Config.SessionToken)
	}

	if p.Config.DryRun {
		p.Config.Publish = false
	} else {
		p.Config.Publish = true
	}

	input := &lambda.UpdateFunctionCodeInput{}
	input.SetDryRun(p.Config.DryRun)
	input.SetFunctionName(p.Config.FunctionName)
	input.SetPublish(p.Config.Publish)

	if p.Config.ImageURI != "" {
		input.SetImageUri(p.Config.ImageURI)
	}

	if p.Config.ReversionID != "" {
		input.SetRevisionId(p.Config.ReversionID)
	}

	if p.Config.S3Bucket != "" && p.Config.S3Key != "" {
		input.SetS3Key(p.Config.S3Key)
		input.SetS3Bucket(p.Config.S3Bucket)

		if p.Config.S3ObjectVersion != "" {
			input.SetS3ObjectVersion(p.Config.S3ObjectVersion)
		}
	}

	sources := trimValues(p.Config.Source)
	if len(sources) != 0 {
		files := globList(sources)
		path := os.TempDir() + "/output.zip"
		zip := archiver.NewZip()
		if len(files) != 0 {
			if err := zip.Archive(files, path); err != nil {
				return err
			}

			p.Config.ZipFile = path
		}
	}

	if p.Config.ZipFile != "" {
		contents, err := os.ReadFile(p.Config.ZipFile)
		if err != nil {
			return err
		}

		input.SetZipFile(contents)
	}

	isUpdateConfig := false
	cfg := &lambda.UpdateFunctionConfigurationInput{}
	cfg.SetFunctionName(p.Config.FunctionName)
	if p.Config.MemorySize > 0 {
		isUpdateConfig = true
		cfg.SetMemorySize(p.Config.MemorySize)
	}
	if p.Config.Timeout > 0 {
		isUpdateConfig = true
		cfg.SetTimeout(p.Config.Timeout)
	}
	if len(p.Config.Handler) > 0 {
		isUpdateConfig = true
		cfg.SetHandler(p.Config.Handler)
	}
	if len(p.Config.Role) > 0 {
		isUpdateConfig = true
		cfg.SetRole(p.Config.Role)
	}
	if len(p.Config.Runtime) > 0 {
		isUpdateConfig = true
		cfg.SetRuntime(p.Config.Runtime)
	}
	if p.Config.Description != "" {
		isUpdateConfig = true
		cfg.SetDescription(p.Config.Description)
	}
	if len(p.Config.Layers) > 0 {
		isUpdateConfig = true
		var layers []*string
		for _, v := range trimValues(p.Config.Layers) {
			layers = append(layers, aws.String(v))
		}
		cfg.SetLayers(layers)
	}

	envs := trimValues(p.Config.Environment)
	if len(envs) > 0 {
		isUpdateConfig = true
		cfg.SetEnvironment(p.loadEnvironment(envs))
	}

	subnets := trimValues(p.Config.Subnets)
	securityGroups := trimValues(p.Config.SecurityGroups)
	if len(subnets) > 0 || len(securityGroups) > 0 {
		isUpdateConfig = true
		cfg.SetVpcConfig(&lambda.VpcConfig{
			SubnetIds:        aws.StringSlice(subnets),
			SecurityGroupIds: aws.StringSlice(securityGroups),
		})
	}

	if p.Config.TracingMode != "" {
		isUpdateConfig = true
		cfg.SetTracingConfig(&lambda.TracingConfig{
			Mode: aws.String(p.Config.TracingMode),
		})
	}

	svc := lambda.New(sess, config)

	if isUpdateConfig {
		if err := svc.WaitUntilFunctionUpdatedV2WithContext(
			aws.BackgroundContext(),
			&lambda.GetFunctionInput{
				FunctionName: aws.String(p.Config.FunctionName),
			},
			request.WithWaiterMaxAttempts(p.Config.MaxAttempts),
		); err != nil {
			log.Println(err.Error())
			return err
		}

		// UpdateFunctionConfiguration API operation for AWS Lambda.
		result, err := svc.UpdateFunctionConfiguration(cfg)
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
				case lambda.ErrCodeResourceConflictException:
					log.Println(lambda.ErrCodeResourceConflictException, aerr.Error())
				case lambda.ErrCodePreconditionFailedException:
					log.Println(lambda.ErrCodePreconditionFailedException, aerr.Error())
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

		p.dump(result)
	}

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

	p.dump(result)

	return nil
}

func (p *Plugin) dump(val ...any) {
	if !p.Config.Debug {
		return
	}

	dump.P(val)
}

func trimValues(keys []string) []string {
	var newKeys []string

	for _, value := range keys {
		value = strings.TrimSpace(value)
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
