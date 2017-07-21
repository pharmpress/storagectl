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
	app.Usage = "A simple command line client for azure storage."
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "account-name", Value: "", Usage: "Account name"},
		cli.StringFlag{Name: "account-key", Value: "", Usage: "Account key"},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name: "azure",
			Subcommands: []cli.Command{
				command.NewS3LsCommand(),
				command.NewS3UploadCommand(),
				command.NewS3DownloadCommand(),
			},
		},
		cli.Command{
			Name: "s3",
			Subcommands: []cli.Command{
				command.NewAzureLsCommand(),
				command.NewAzureUploadCommand(),
				command.NewAzureDownloadCommand(),
			},
		},
	}

	app.Run(os.Args)
	fmt.Println("Finish.")
}
