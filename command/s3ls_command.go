package command

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pharmpress/storagectl/util"
	"github.com/urfave/cli"
)

// NewS3LsCommand is ls implementation for s3
func NewS3LsCommand() cli.Command {
	return cli.Command{
		Name:  "ls",
		Usage: "retrieve container blobs",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "container", Usage: "Container Name"},
		},
		Action: func(c *cli.Context) {
			s3lsCommandFunc(c)
		},
	}
}

// lsCommandFunc executes the "ls" command.
func s3lsCommandFunc(c *cli.Context) error {
	accountName := c.GlobalString("account-name")
	accountKey := c.GlobalString("account-key")
	containerName := c.String("container")

	s3Svc, err := util.NewS3ServiceInstance(accountName, accountKey, containerName)
	if err != nil {
		return err
	}
	resp, err := s3Svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(containerName),
	})
	if err != nil {
		return err
	}
	for _, key := range resp.Contents {
		fmt.Println(fmt.Sprintf("blob -> %+v", *key.Key))
	}

	return err
}
