package util

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// NewS3ServiceInstance centralises s3client creation
func NewS3ServiceInstance(accountName string, accountKey string, containerName string) (*s3.S3, error) {
	token := ""
	creds := credentials.NewStaticCredentials(accountName, accountKey, token)
	_, err := creds.Get()
	if err != nil {
		fmt.Printf("bad credentials: %s", err)
		return nil, err
	}

	region, err := s3manager.GetBucketRegion(aws.BackgroundContext(), session.Must(session.NewSession()), containerName, "eu-west-1")
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			fmt.Fprintf(os.Stderr, "unable to find bucket %s's region not found\n", containerName)
		}
		return nil, err
	}

	cfg := aws.NewConfig().WithCredentials(creds).WithRegion(region)
	s3Svc := s3.New(session.Must(session.NewSession()), cfg)
	return s3Svc, err
}
