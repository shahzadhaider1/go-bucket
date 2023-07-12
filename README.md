# go-bucket

A small tool to interact with your S3-compatible Cloud Object Storage and perform various operations. 

## Usage

Here is how you can use this in your code:

1. Create a new file `main.go` and paste the following code in it
2. Run `go mod tidy`
3. Run `go run main.go`

```
package main

import (
	"fmt"

	"github.com/shahzadhaider1/go-bucket/bucket"
)

func main() {
	creds := &bucket.Credentials{
		AccessKey:  "access_key",
		SecretKey:  "secret_key",
		Endpoint:   "endpoint",
		BucketName: "bucket_name",
	}

	if err := creds.ClearBucket(); err != nil {
		fmt.Println("failed to clear bucket, err : ", err)
	}

}
```

