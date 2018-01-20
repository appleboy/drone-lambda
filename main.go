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
	}

	app.Version = Version

	if BuildNum != "" {
		app.Version = app.Version + "+" + BuildNum
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Println(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Config: Config{
			Region:    c.String("region"),
			AccessKey: c.String("access-key"),
			SecretKey: c.String("secret-key"),
		},
	}

	return plugin.Exec()
}
