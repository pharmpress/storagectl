package main

import (
	"fmt"
	"os"

	"github.com/pharmpress/storagectl/command"
	"github.com/pharmpress/storagectl/version"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "storagectl"
	app.Version = version.Version
	app.Usage = "A simple command line client for azure and s3 storage."
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "account-name", Value: "", Usage: "Account name"},
		cli.StringFlag{Name: "account-key", Value: "", Usage: "Account key"},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:  "s3",
			Usage: "options for s3 storage",
			Subcommands: []cli.Command{
				command.NewS3LsCommand(),
				command.NewS3UploadCommand(),
				command.NewS3DownloadCommand(),
			},
		},
		cli.Command{
			Name:  "azure",
			Usage: "options for azure storage",
			Subcommands: []cli.Command{
				command.NewAzureLsCommand(),
				command.NewAzureUploadCommand(),
				command.NewAzureDownloadCommand(),
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Finish.")
	}
}
