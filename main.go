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
			EnvVar: "PLUGIN_REGION,AWS_REGION,INPUT_AWS_REGION",
			Value:  "us-east-1",
		},
		cli.StringFlag{
			Name:   "access-key",
			Usage:  "AWS ACCESS KEY",
			EnvVar: "PLUGIN_ACCESS_KEY,AWS_ACCESS_KEY_ID,INPUT_AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "secret-key",
			Usage:  "AWS SECRET KEY",
			EnvVar: "PLUGIN_SECRET_KEY,AWS_SECRET_ACCESS_KEY,INPUT_AWS_SECRET_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "session-token",
			Usage:  "AWS Session token",
			EnvVar: "PLUGIN_SESSION_TOKEN,AWS_SESSION_TOKEN,INPUT_AWS_SESSION_TOKEN",
		},
		cli.StringFlag{
			Name:   "aws-profile",
			Usage:  "AWS profile",
			EnvVar: "PLUGIN_PROFILE,AWS_PROFILE,INPUT_AWS_PROFILE",
		},
		cli.StringFlag{
			Name:   "function-name",
			Usage:  "AWS lambda function name",
			EnvVar: "PLUGIN_FUNCTION_NAME,FUNCTION_NAME,INPUT_FUNCTION_NAME",
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
			Source:          c.StringSlice("source"),
			DryRun:          c.Bool("dry-run"),
		},
	}

	return plugin.Exec()
}
