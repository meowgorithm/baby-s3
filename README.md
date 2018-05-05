Baby S3
========

AWS is really complicated and my needs are usually not that complicated. This
library exists to handle some common S3 operations in a simple way.

For example:

    import (
        "github.com/meowgorithm/baby-s3"
    )

    func main() {
        if err := s3.CreateBucket("my-cute-bucket"); err != nil  {
            fmt.Println("It didn't work :(", err)
            return
        }

        if err := s3.MakeBucketPublic("my-cute-bucket"); err != nil {
            fmt.Println("It didn't work. Weird :/", err)
        }

        err := s3.UploadObject("my-cute-bucket", "some-bytes.txt", []byte("a few nice bytes"));
        err != nil {
            fmt.Println("Nope :(", err)
        }
    }


⚠️ For now, this library is in an alpha state and the API could change. If you
have any thoughts about the API do let me know.

## AWS Keys and Regions and Stuff

AWS loves it when you put stuff in environment variables, so you'll need to do
that to work with this library. Set the following:

    AWS_ACCESS_KEY
    AWS_SECRET_ACCESS_KEY
    AWS_REGION
