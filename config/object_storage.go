package config

import (
    "os"
    "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

var S3 *s3.S3

// func InitializeS3() {
//     session, err := session.NewSession(&aws.Config{
//         Region: aws.String(os.Getenv("XYZ1_REGION")),
//         Endpoint: aws.String(os.Getenv("XYZ1_ENDPOINT")),
//         Credentials: credentials.NewStaticCredentials(
//             os.Getenv("XYZ1_ACCESS"),
//             os.Getenv("XYZ1_SECRET"),
//             "",
//         ),
//     })

//     if err != nil {
//         log.Fatal(err)
//     }

//     S3 = s3.New(session)
// }

func InitializeS3() {
    cfg := &aws.Config{
        Region: aws.String(os.Getenv("XYZ1_REGION")),
        Endpoint: aws.String(os.Getenv("XYZ1_ENDPOINT")),
        Credentials: credentials.NewStaticCredentials(
            os.Getenv("XYZ1_ACCESS"),
            os.Getenv("XYZ1_SECRET"),
            "",
        ),
    }

    sess := session.Must(session.NewSession(cfg))
    S3 = s3.New(sess)
}