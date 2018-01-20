package main

import (
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Version set at compile-time
var (
	Version  string
	BuildNum string
)

func main() {
	app := cli.NewApp()
	app.Name = "Drone Lambda"
	app.Usage = "Deploying Lambda code with drone CI to an existing function"
	app.Copyright = "Copyright (c) 2018 Bo-Yi Wu"
	app.Authors = []cli.Author{
		{
			Name:  "Bo-Yi Wu",
			Email: "appleboy.tw@gmail.com",
		},
	}
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "region",
			Usage:  "AWS Region",
			EnvVar: "PLUGIN_REGION,AWS_REGION",
			Value:  "us-east-1",
		},
		cli.StringFlag{
			Name:   "access-key",
			Usage:  "AWS ACCESS KEY",
			EnvVar: "PLUGIN_ACCESS_KEY,AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "secret-key",
			Usage:  "AWS SECRET KEY",
			EnvVar: "PLUGIN_SECRET_KEY,AWS_SECRET_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "session-token",
			Usage:  "AWS Session token",
			EnvVar: "PLUGIN_SESSION_TOKEN,AWS_SESSION_TOKEN",
		},
		cli.StringFlag{
			Name:   "aws-profile",
			Usage:  "AWS profile",
			EnvVar: "PLUGIN_PROFILE,AWS_PROFILE",
		},
		cli.StringFlag{
			Name:   "function-name",
			Usage:  "AWS lambda function name",
			EnvVar: "PLUGIN_FUNCTION_NAME",
		},
		cli.StringFlag{
			Name:   "s3-bucket",
			Usage:  "AWS lambda S3 bucket",
			EnvVar: "PLUGIN_S3_BUCKET",
		},
		cli.StringFlag{
			Name:   "s3-key",
			Usage:  "AWS lambda S3 bucket key",
			EnvVar: "PLUGIN_S3_KEY",
		},
		cli.StringFlag{
			Name:   "s3-object-version",
			Usage:  "AWS lambda s3 object version",
			EnvVar: "PLUGIN_S3_OBJECT_VERSION",
		},
		cli.StringFlag{
			Name:   "zip-file",
			Usage:  "AWS lambda zip file",
			EnvVar: "PLUGIN_ZIP_FILE",
		},
	}

	app.Version = Version

	if BuildNum != "" {
		app.Version = app.Version + "+" + BuildNum
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Warningln(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

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
		},
	}

	return plugin.Exec()
}
