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
			Name:   "secret, s",
			Usage:  "line channel secret",
			EnvVar: "PLUGIN_CHANNEL_SECRET,LINE_CHANNEL_SECRET",
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
			Region: c.String("secret"),
		},
	}

	return plugin.Exec()
}
