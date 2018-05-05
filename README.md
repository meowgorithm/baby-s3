Baby AWS
========

AWS is really complicated and my needs are usually not that complicated. Thus,
I've built this library to handle some common AWS operations in a simple way.

For example:

    import (
        "github.com/meowgorithm/baby-aws/s3"
    )

    func main() {
        if err := s3.CreateBucket("my-cute-bucket"); err != nil  {
            fmt.Println("It didn't work :(", err)
            return
        }

        err := s3.UploadObject("my-cute-bucket", "some-bytes.txt", []byte("a few nice bytes"));
        err != nil {
            fmt.Println("Nope :(", err)
        }
    }


⚠️ For now, this library is in an alpha state and the API may change. If you
have any thoughts about the API let me know.

## AWS Keys and Regions and Stuff

AWS loves it when you put stuff in environment variables, so you'll need to do
that to work with this library. Set the following variables:

    AWS_ACCESS_KEY
    AWS_SECRET_ACCESS_KEY
    AWS_REGION
