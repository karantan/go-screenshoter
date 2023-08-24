# go-screenshoter

`go-screenshoter` is a simple yet effective project that demonstrates the integration of
several powerful tools to capture screenshots of web pages and upload them to Cloudflare's R2
storage. Specifically crafted as a hands-on example, it showcases the prowess of Go
language, AWS Lambda, Serverless framework, and the flexibility of `devenv` for reproducible
developer environments.

## Tools Used

- **devenv (devenv.sh)**: Enables fast, declarative, reproducible, and composable developer
environments using Nix.
- **Go Language**: The backbone of our application, providing efficiency and concurrency.
- **AWS Lambda**: Our serverless compute service where our function resides and gets executed.
- **Serverless Framework**: Facilitates deploying and managing applications on cloud platforms
without worrying about infrastructure.

## Prerequisites

Before diving into `go-screenshoter`, ensure you have the following installed:

- **Nix Language**: Essential for our `devenv` tool.
- **devenv tool**: Install it by following guidelines [here](devenv.sh).

## Getting Started

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/karantan/go-screenshoter.git
   cd go-screenshoter
   ```

2. **Setup Development Environment**:

   With `devenv` and `Nix` already installed:

   ```bash
   direnv shell
   ```

   This will setup a consistent and reproducible developer environment for you.

3. **Deploy to AWS Lambda**:

   Next you'll need to download the chromium binary compatible with x86_64.
   E.g. [alixaxel/chrome-aws-lambda](https://raw.githubusercontent.com/alixaxel/chrome-aws-lambda/master/bin/chromium.br)

   Extract it in the `layer` folder and make sure it has executable permissions
   (`chmod 777 chromium`).

   Follow Serverless framework guidelines to deploy the function to AWS Lambda. Ensure
   your AWS credentials are properly set up.

   ```bash
   make deploy
   ```

4. **Usage**:

   Once deployed, invoke your Lambda function with a target URL to get a presigned URL
   of the captured screenshot.

   Example:

   ```bash
   sls invoke --function screenshot --path=lib/data.json
   ```

## Project Structure and Insights

### Lambda Layer for Chromium

To enable browser capabilities within the AWS Lambda environment, a layer containing the
Chromium binary was necessary. Here's how it's set up in the `serverless.yml` configuration:

```yml
layers:
  chromium:
    path: layer
    package:
      include:
        - ./**

functions:
  screenshot:
    handler: bin/screenshoter
    layers:
      - {Ref: ChromiumLambdaLayer} # ${CamelCaseLayerName} + LambdaLayer
```

The Chromium binary resides within the `layer` folder of this repository. It was
downloaded from [alixaxel/chrome-aws-lambda](https://raw.githubusercontent.com/alixaxel/chrome-aws-lambda/master/bin/chromium.br)

### Helpful Documentation

During the development of this project, certain resources proved invaluable:

- [Serverless Framework on AWS Lambda Layers](https://www.serverless.com/framework/docs/providers/aws/guide/layers/)
- [Article on using AWS Lambda Layers with the Serverless Framework](https://www.serverless.com/blog/publish-aws-lambda-layers-serverless-framework/)
- [AWS Lambda Layers documentation](https://docs.aws.amazon.com/lambda/latest/dg/chapter-layers.html)

### Resolving File Permissions in Lambda

The files from Lambda layers get extracted to the `/opt` directory. Initially, the
Chromium executable in this directory was not set as executable (missing chmod 777).

To diagnose this, the following function was utilized to walk and list all files in the `/opt` directory:

```go
func walk(path string) {
    filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }
        fmt.Println(path, info.Size())
        return nil
    })
}
```

This provided the insight that Chromium was present but lacked execution permissions.

### Browser Control with Go

For browser control, the [go-rod](https://go-rod.github.io/) library was adopted.
A particularly helpful [repository](https://github.com/YoungiiJC/go-rod-aws-lambda/) provided
guidelines on configuring browser settings specifically for the AWS Lambda environment.

## Contributions

Feel free to fork, raise issues, or submit Pull Requests. All contributions are welcome!
