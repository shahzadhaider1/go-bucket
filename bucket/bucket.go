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
	AccessKey  string
	SecretKey  string
	Endpoint   string
	BucketName string
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
		Endpoint:         aws.String(c.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
	})
	if err != nil {
		return fmt.Errorf("failed to create session, error : %v", err)
	}

	// Create an S3 service client
	svc := s3.New(sess)

	for {
		// List objects in the bucket
		listInput := &s3.ListObjectsInput{
			Bucket: aws.String(c.BucketName),
		}

		listOutput, err := svc.ListObjects(listInput)
		if err != nil {
			fmt.Errorf("failed to list the objects, error : %v", err)
			return err
		}

		if len(listOutput.Contents) == 0 {
			// Bucket is empty, exit the loop
			break
		}

		// Delete each object in the bucket
		for _, obj := range listOutput.Contents {
			deleteInput := &s3.DeleteObjectInput{
				Bucket: aws.String(c.BucketName),
				Key:    obj.Key,
			}
			if _, err = svc.DeleteObject(deleteInput); err != nil {
				return fmt.Errorf("failed to delete the object, key : %v - error : %v", *obj.Key, err)
			}
		}
	}

	return nil
}
