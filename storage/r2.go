// See https://developers.cloudflare.com/r2/examples/aws/aws-sdk-go/
package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func client() *s3.Client {
	accountId := os.Getenv("CF_ACCOUNT_ID")
	accessKeyId := os.Getenv("CF_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("CF_ACCESS_KEY_SECRET")

	if accountId == "" {
		log.Fatalf("Missing CF_ACCOUNT_ID env. var.")
	}
	if accessKeyId == "" {
		log.Fatalf("Missing CF_ACCESS_KEY_ID env. var.")
	}
	if accessKeySecret == "" {
		log.Fatalf("Missing CF_ACCESS_KEY_SECRET env. var.")
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	return s3.NewFromConfig(cfg)

}

func Upload(objectKey, fileName string) error {
	bucketName := os.Getenv("CF_BUCKET_NAME")

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	c := client()
	_, err = c.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	return err
}

func GetPresignURL(objectKey string) (string, error) {
	bucketName := os.Getenv("CF_BUCKET_NAME")
	c := client()
	presignClient := s3.NewPresignClient(c)
	presignedUrl, err := presignClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		},
		s3.WithPresignExpires(time.Minute*15))
	if err != nil {
		return "", err
	}
	return presignedUrl.URL, nil
}
