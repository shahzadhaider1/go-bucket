# go-bucket

A small tool to interact with your S3-compatible Cloud Object Storage and perform various operations. 

Here is how you can use this in your code: 

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
		fmt.Errorf("failed to clear bucket, err %v: ", err)
	}

}
```

