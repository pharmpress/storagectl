package command

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pharmpress/storagectl/util"
	"github.com/urfave/cli"
)

// NewS3DownloadCommand is download implementation for s3
func NewS3DownloadCommand() cli.Command {
	return cli.Command{
		Name:  "download",
		Usage: "download blob from a container",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "container", Usage: "Container Name"},
			cli.StringFlag{Name: "blob", Usage: "Blob name"},
		},
		Action: func(c *cli.Context) {
			s3DownloadCommandFunc(c)
		},
	}
}

func s3DownloadCommandFunc(c *cli.Context) error {
	accountName := c.GlobalString("account-name")
	accountKey := c.GlobalString("account-key")
	containerName := c.String("container")
	blobName := c.String("blob")
	fileToDownload := "/file_to_download"

	if len(c.Args()) != 0 {
		fileToDownload = c.Args()[0]
	}

	s3Svc, err := util.NewS3ServiceInstance(accountName, accountKey, containerName)
	if err != nil {
		return err
	}

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloaderWithClient(s3Svc)

	// Create a file to write the S3 Object contents to.
	f, err := os.Create(fileToDownload)
	if err != nil {
		return fmt.Errorf("failed to create file %q, %v", fileToDownload, err)
	}

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(containerName),
		Key:    aws.String(blobName),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file downloaded, %d bytes\n", n)

	return err
}
