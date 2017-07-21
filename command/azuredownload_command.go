package command

import (
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/urfave/cli"
)

// NewAzureDownloadCommand is download implementation for Azure
func NewAzureDownloadCommand() cli.Command {
	return cli.Command{
		Name:  "download",
		Usage: "download blob from a container",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "container", Usage: "Container Name"},
			cli.StringFlag{Name: "blob", Usage: "Blob name"},
		},
		Action: func(c *cli.Context) {
			azureDownloadCommandFunc(c)
		},
	}
}

func azureDownloadCommandFunc(c *cli.Context) error {
	accountName := c.GlobalString("account-name")
	accountKey := c.GlobalString("account-key")
	containerName := c.String("container")
	blobName := c.String("blob")
	fileToDownload := "/file_to_download"

	if len(c.Args()) != 0 {
		fileToDownload = c.Args()[0]
	}

	client, err := storage.NewBasicClient(accountName, accountKey)

	if err != nil {
		fmt.Println(err)
		return err
	}

	blobService := client.GetBlobService()

	file, err := os.Create(fileToDownload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	fmt.Println("Downloading " + fileToDownload)

	blob, err := blobService.GetBlob(containerName, blobName)
	if err != nil {
		fmt.Println(err)
	}
	defer blob.Close()

	_, err = io.Copy(file, blob)

	if err != nil {
		fmt.Println(err)
	}

	return err
}
