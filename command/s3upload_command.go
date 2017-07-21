package command

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pharmpress/storagectl/util"
	"github.com/urfave/cli"
)

// NewS3UploadCommand is upload implementation for s3
func NewS3UploadCommand() cli.Command {
	return cli.Command{
		Name:  "upload",
		Usage: "upload blob into a container",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "container", Usage: "Container Name"},
			cli.StringFlag{Name: "blob", Usage: "Blob name"},
		},
		Action: func(c *cli.Context) {
			s3UploadCommandFunc(c)
		},
	}
}

func s3UploadCommandFunc(c *cli.Context) error {
	accountName := c.GlobalString("account-name")
	accountKey := c.GlobalString("account-key")
	containerName := c.String("container")
	blobName := c.String("blob")
	fileToUpload := "/file_to_upload"

	if len(c.Args()) != 0 {
		fileToUpload = c.Args()[0]
	}

	s3Svc, err := util.NewS3ServiceInstance(accountName, accountKey, containerName)
	if err != nil {
		return err
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploaderWithClient(s3Svc)

	f, err := os.Open(fileToUpload)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", fileToUpload, err)
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(containerName),
		Key:    aws.String(blobName),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)

	return err
}
