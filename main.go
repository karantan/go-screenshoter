// Package main provides functionality to capture screenshots of web pages and then uploads them
// to Cloudflare's R2 storage system, specifically designed for deployment on AWS Lambda.
//
// The main workflow of this package is as follows:
// 1. On invocation, navigate to a specific URL and take a screenshot using the TakeScreenshot function.
// 2. Upload the taken screenshot to Cloudflare's R2 storage using the Upload function.
// 3. Retrieve a presigned URL for the uploaded screenshot to allow for easy viewing and sharing.
// 4. Return the presigned URL as the Lambda function's response.
//
// Note: For the package to function properly in the AWS Lambda environment:
// - Ensure that the required browser drivers are included in the Lambda deployment package.
// - Set up necessary AWS Lambda environment variables or permissions for Cloudflare's R2 storage access.
// - Consider any timeout restrictions of Lambda when navigating to websites and capturing screenshots.
//
// The Lambda function expects an event input with details like the target URL, and it responds with
// the presigned URL of the stored screenshot.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path"

	"screenshoter/logger"
	"screenshoter/naviga"
	"screenshoter/storage"
	"screenshoter/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/joho/godotenv/autoload"
)

var log = logger.New("main", false)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Body is a struct for holding request post body content
type Body struct {
	URL string `json:"url"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	body := Body{}
	json.Unmarshal([]byte(request.Body), &body)

	if body.URL == "" {
		return Response{StatusCode: 400}, nil
	}

	hostname, err := utils.GetDomainFromURL(body.URL)
	if err != nil {
		return Response{StatusCode: 503}, err
	}

	screenshotPath := path.Join(os.Getenv("TEMPORARY_STORAGE"), hostname+".png")

	err = naviga.TakeScreenshot(body.URL, screenshotPath)
	if err != nil {
		return Response{StatusCode: 503}, err
	}

	err = storage.Upload(hostname, screenshotPath)
	if err != nil {
		return Response{StatusCode: 503}, err
	}

	presignURL, err := storage.GetPresignURL(hostname)
	if err != nil {
		return Response{StatusCode: 503}, err
	}

	respBody, err := json.Marshal(map[string]interface{}{
		"url":        body.URL,
		"screenshot": presignURL,
	})
	if err != nil {
		return Response{StatusCode: 503}, err
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, respBody)
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "dev" {
		url := flag.String("url", "https://www.wikipedia.org/", "Url of a website you want to make a screenshot")
		screenshotPath := flag.String("path", "lib/wikipedia.org.png", "Path where the screenshot will be stored")
		flag.Parse()

		err := naviga.TakeScreenshot(*url, *screenshotPath)
		if err != nil {
			log.Error(err)
			return
		}
		log.Infof("Screenshot successfuly created: %s", screenshotPath)

	} else {
		lambda.Start(Handler)
	}
}
