package command

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/urfave/cli"
)

// NewAzureLsCommand is ls implementation for Azure
func NewAzureLsCommand() cli.Command {
	return cli.Command{
		Name:  "ls",
		Usage: "retrieve container blobs",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "container", Usage: "Container Name"},
		},
		Action: func(c *cli.Context) {
			azureLsCommandFunc(c)
		},
	}
}

// lsCommandFunc executes the "ls" command.
func azureLsCommandFunc(c *cli.Context) error {
	accountName := c.GlobalString("account-name")
	accountKey := c.GlobalString("account-key")
	containerName := c.String("container")
	client, err := storage.NewBasicClient(accountName, accountKey)

	if err != nil {
		fmt.Println(err)
		return err
	}

	blobService := client.GetBlobService()

	blobs, err := blobService.ListBlobs(containerName, storage.ListBlobsParameters{})

	if err != nil {
		fmt.Println(err)
	} else {
		for _, file := range blobs.Blobs {
			fmt.Println(fmt.Sprintf("blob -> %+v", file))
		}
	}
	return err
}
