package bucket

import (
	"fmt"
	"sync"

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
	Region     string
}

// NewS3Client initialises a new S3 client
func (c *Credentials) NewS3Client(accessKey, secretKey, endpoint, bucketName, region string) *Credentials {
	return &Credentials{
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		Endpoint:   endpoint,
		BucketName: bucketName,
		Region:     region,
	}
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
		Region:           aws.String(c.Region), // Specify the appropriate region
		Endpoint:         aws.String(c.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
	})
	if err != nil {
		return fmt.Errorf("failed to create session, error : %v", err)
	}

	// Create an S3 service client
	svc := s3.New(sess)

	// Create a channel to receive deletion errors
	errCh := make(chan error)

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

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
			wg.Add(1)
			go func(key *string) {
				defer wg.Done()

				deleteInput := &s3.DeleteObjectInput{
					Bucket: aws.String(c.BucketName),
					Key:    key,
				}
				_, err := svc.DeleteObject(deleteInput)
				if err != nil {
					errCh <- fmt.Errorf("failed to delete the object, key : %v - error : %v", *key, err)
				}
			}(obj.Key)
		}

		// Wait for all deletions to complete
		wg.Wait()

		// Check if there were any deletion errors
		select {
		case err := <-errCh:
			return err
		default:
			// No errors, continue with the next batch of objects
		}
	}

	return nil
}

// ListObjects lists all the objects available in the COS Bucket
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
//	An error if the listing process encounters any issues, or nil if successful.
func (c *Credentials) ListObjects() error {
	// Create a session
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(c.Region), // Replace with your region
		Endpoint:         aws.String(c.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
	})
	if err != nil {
		panic(err)
	}

	// Create an S3 client
	svc := s3.New(sess)

	// Create a channel to receive deletion errors
	errCh := make(chan error)

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	for {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// List objects in the bucket
			resp, err := svc.ListObjects(&s3.ListObjectsInput{
				Bucket: aws.String(c.BucketName),
			})

			if err != nil {
				panic(err)
			}

			// Iterate through the objects and print their names
			for _, obj := range resp.Contents {
				fmt.Println("Object Name:", *obj.Key)
			}
		}()
	}

	// Wait for all deletions to complete
	wg.Wait()

	// Check if there were any deletion errors
	select {
	case err = <-errCh:
		return err
	default:
		// No errors, continue with the next batch of objects
	}

	return nil
}
