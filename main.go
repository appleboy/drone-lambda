package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

// Version set at compile-time
var (
	Version string
)

func main() {
	// Load env-file if it exists first
	if filename, found := os.LookupEnv("PLUGIN_ENV_FILE"); found {
		_ = godotenv.Load(filename)
	}

	app := cli.NewApp()
	app.Name = "Drone Lambda"
	app.Usage = "Deploying Lambda code with drone CI to an existing function"
	app.Copyright = "Copyright (c) 2020 Bo-Yi Wu"
	app.Authors = []cli.Author{
		{
			Name:  "Bo-Yi Wu",
			Email: "appleboy.tw@gmail.com",
		},
	}
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "region",
			Usage:  "AWS Region",
			EnvVar: "PLUGIN_REGION,PLUGIN_AWS_REGION,AWS_REGION,INPUT_AWS_REGION",
			Value:  "us-east-1",
		},
		cli.StringFlag{
			Name:   "access-key",
			Usage:  "AWS ACCESS KEY",
			EnvVar: "PLUGIN_ACCESS_KEY,PLUGIN_AWS_ACCESS_KEY_ID,AWS_ACCESS_KEY_ID,INPUT_AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "secret-key",
			Usage:  "AWS SECRET KEY",
			EnvVar: "PLUGIN_SECRET_KEY,PLUGIN_AWS_SECRET_ACCESS_KEY,AWS_SECRET_ACCESS_KEY,INPUT_AWS_SECRET_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "session-token",
			Usage:  "AWS Session token",
			EnvVar: "PLUGIN_SESSION_TOKEN,PLUGIN_AWS_SESSION_TOKEN,AWS_SESSION_TOKEN,INPUT_AWS_SESSION_TOKEN",
		},
		cli.StringFlag{
			Name:   "aws-profile",
			Usage:  "AWS profile",
			EnvVar: "PLUGIN_PROFILE,PLUGIN_AWS_PROFILE,AWS_PROFILE,INPUT_AWS_PROFILE",
		},
		cli.StringFlag{
			Name:   "function-name",
			Usage:  "AWS lambda function name",
			EnvVar: "PLUGIN_FUNCTION_NAME,FUNCTION_NAME,INPUT_FUNCTION_NAME",
		},
		cli.StringFlag{
			Name:   "reversion-id",
			Usage:  "Only update the function if the revision ID matches the ID that's specified.",
			EnvVar: "PLUGIN_REVERSION_ID,REVERSION_ID,INPUT_REVERSION_ID",
		},
		cli.StringFlag{
			Name:   "s3-bucket",
			Usage:  "An Amazon S3 bucket in the same AWS Region as your function. The bucket can be in a different AWS account.",
			EnvVar: "PLUGIN_S3_BUCKET,S3_BUCKET,INPUT_S3_BUCKET",
		},
		cli.StringFlag{
			Name:   "s3-key",
			Usage:  "The Amazon S3 key of the deployment package.",
			EnvVar: "PLUGIN_S3_KEY,S3_KEY,INPUT_S3_KEY",
		},
		cli.StringFlag{
			Name:   "s3-object-version",
			Usage:  "AWS lambda s3 object version",
			EnvVar: "PLUGIN_S3_OBJECT_VERSION,S3_OBJECT_VERSION,INPUT_S3_OBJECT_VERSION",
		},
		cli.StringFlag{
			Name:   "zip-file",
			Usage:  "AWS lambda zip file",
			EnvVar: "PLUGIN_ZIP_FILE,ZIP_FILE,INPUT_ZIP_FILE",
		},
		cli.StringSliceFlag{
			Name:   "source",
			Usage:  "zip file list",
			EnvVar: "PLUGIN_SOURCE,SOURCE,INPUT_SOURCE",
		},
		cli.BoolFlag{
			Name:   "dry-run",
			Usage:  "Set to true to validate the request parameters and access permissions without modifying the function code.",
			EnvVar: "PLUGIN_DRY_RUN,DRY_RUN,INPUT_DRY_RUN",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "Show debug message after upload the lambda successfully.",
			EnvVar: "PLUGIN_DEBUG,DEBUG,INPUT_DEBUG",
		},
		cli.BoolTFlag{
			Name:   "publish",
			Usage:  "Set to true to publish a new version of the function after updating the code.",
			EnvVar: "PLUGIN_PUBLISH,PUBLISH,INPUT_PUBLISH",
		},
		cli.Int64Flag{
			Name:   "memory-size",
			Usage:  "The amount of memory that your function has access to. Increasing the function's memory also increases its CPU allocation. The default value is 128 MB. The value must be a multiple of 64 MB.",
			EnvVar: "PLUGIN_MEMORY_SIZE,MEMORY_SIZE,INPUT_MEMORY_SIZE",
		},
		cli.Int64Flag{
			Name:   "timeout",
			Usage:  "The amount of time that Lambda allows a function to run before stopping it. The default is 3 seconds. The maximum allowed value is 900 seconds.",
			EnvVar: "PLUGIN_TIMEOUT,TIMEOUT,INPUT_TIMEOUT",
		},
		cli.StringFlag{
			Name:   "handler",
			Usage:  "The name of the method within your code that Lambda calls to execute your function.",
			EnvVar: "PLUGIN_HANDLER,HANDLER,INPUT_HANDLER",
		},
		cli.StringFlag{
			Name:   "role",
			Usage:  "The Amazon Resource Name (ARN) of the function's execution role.",
			EnvVar: "PLUGIN_ROLE,ROLE,INPUT_ROLE",
		},
		cli.StringFlag{
			Name:   "runtime",
			Usage:  "The identifier of the function's runtime.",
			EnvVar: "PLUGIN_RUNTIME,RUNTIME,INPUT_RUNTIME",
		},
		cli.StringSliceFlag{
			Name:   "environment",
			Usage:  "Lambda Environment variables",
			EnvVar: "PLUGIN_ENVIRONMENT,ENVIRONMENT,INPUT_ENVIRONMENT",
		},
		cli.StringSliceFlag{
			Name:   "layers",
			Usage:  "A list of function layers",
			EnvVar: "PLUGIN_LAYERS,LAYERS,INPUT_LAYERS",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA,GITHUB_SHA",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "image-uri",
			Usage:  "URI of a container image in the Amazon ECR registry.",
			EnvVar: "PLUGIN_IMAGE_URI,IMAGE_URI,INPUT_IMAGE_URI",
		},
		cli.StringSliceFlag{
			Name:   "subnets",
			Usage:  "Select the VPC subnets for Lambda to use to set up your VPC configuration.",
			EnvVar: "PLUGIN_SUBNETS,SUBNETS,INPUT_SUBNETS",
		},
		cli.StringSliceFlag{
			Name:   "securitygroups",
			Usage:  "Choose the VPC security groups for Lambda to use to set up your VPC configuration.",
			EnvVar: "PLUGIN_SECURITY_GROUPS,SECURITY_GROUPS,INPUT_SECURITY_GROUPS",
		},
		cli.StringFlag{
			Name:   "description",
			Usage:  "A description of the function.",
			EnvVar: "PLUGIN_DESCRIPTION,DESCRIPTION,INPUT_DESCRIPTION",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Config: Config{
			Region:          c.String("region"),
			AccessKey:       c.String("access-key"),
			SecretKey:       c.String("secret-key"),
			Profile:         c.String("aws-profile"),
			S3Bucket:        c.String("s3-bucket"),
			S3Key:           c.String("s3-key"),
			S3ObjectVersion: c.String("s3-object-version"),
			ZipFile:         c.String("zip-file"),
			FunctionName:    c.String("function-name"),
			ReversionID:     c.String("reversion-id"),
			Source:          c.StringSlice("source"),
			DryRun:          c.Bool("dry-run"),
			Debug:           c.Bool("debug"),
			Publish:         c.Bool("publish"),
			Timeout:         c.Int64("timeout"),
			MemorySize:      c.Int64("memory-size"),
			Handler:         c.String("handler"),
			Role:            c.String("role"),
			Runtime:         c.String("runtime"),
			Environment:     c.StringSlice("environment"),
			Layers:          c.StringSlice("layers"),
			ImageURI:        c.String("image-uri"),
			Subnets:         c.StringSlice("subnets"),
			SecurityGroups:  c.StringSlice("securitygroups"),
			Description:     c.String("description"),
		},
		Commit: Commit{
			Sha:    c.String("commit.sha"),
			Author: c.String("commit.author"),
		},
	}

	return plugin.Exec()
}
