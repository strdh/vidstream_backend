package utils

import (
    "os"
    "io/ioutil"
    "xyzstream/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func ObjUpload(bucket string, filePath string, objKey string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    params := &s3.PutObjectInput{
        Bucket: aws.String(bucket),
        Key: aws.String(objKey),
        Body: file,
    }

    _, err = config.S3.PutObject(params)
    if err != nil {
        return err
    }

    return nil
}

func ObjRead(bucket string, objKey string) ([]uint8, error) {
    params := &s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key: aws.String(objKey),
    }

    response, err := config.S3.GetObject(params)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    data, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }

    return data, nil
}
