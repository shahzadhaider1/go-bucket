package bucket

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Credentials holds the credentials of bucket that needs to be cleared
type Credentials struct {
	accessKey  string
	secretKey  string
	endpoint   string
	bucketName string
}

// ClearBucket will delete the objects from your IBM Cloud COS Bucket
//
// Parameters:
//
//	accessKey: The access key for authentication.
//	secretKey: The secret key for authentication.
//	endpoint: The endpoint URL of the bucket.
//	bucketName: The name of the bucket to clear.
//
// Returns:
//
//	An error if the bucket clearing process encounters any issues, or nil if successful.
func (c *Credentials) ClearBucket() error {
	// Create a session with your IBM Cloud Object Storage credentials
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-geo"), // Specify the appropriate region
		Endpoint:         aws.String(c.endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(c.accessKey, c.secretKey, ""),
	})
	if err != nil {
		return fmt.Errorf("failed to create session, error : %v", err)
	}

	// Create an S3 service client
	svc := s3.New(sess)

	// List objects in the bucket
	listInput := &s3.ListObjectsInput{
		Bucket: aws.String(c.bucketName),
	}

	listOutput, err := svc.ListObjects(listInput)
	if err != nil {
		fmt.Errorf("failed to list the objects, error : %v", err)
		return err
	}

	// Delete each object in the bucket
	for _, obj := range listOutput.Contents {
		deleteInput := &s3.DeleteObjectInput{
			Bucket: aws.String(c.bucketName),
			Key:    obj.Key,
		}
		if _, err = svc.DeleteObject(deleteInput); err != nil {
			return fmt.Errorf("failed to delete the object, key : %v - error : %v", *obj.Key, err)
		}
	}

	return nil
}
