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
   devenv init
   devenv shell
   ```

   This will setup a consistent and reproducible developer environment for you.

3. **Deploy to AWS Lambda**:

   Next you'll need to download the chromium binary compatible with x86_64.
   E.g. [alixaxel/chrome-aws-lambda](https://raw.githubusercontent.com/alixaxel/chrome-aws-lambda/master/bin/chromium.br)

   Extract it in the `layer` folder and make sure it has executable permissions.

   ```bash
   wget -P layer https://raw.githubusercontent.com/alixaxel/chrome-aws-lambda/master/bin/chromium.br
   brotli --decompress --rm --output=layer/chromium layer/chromium.br
   chmod 777 layer/chromium
   ```

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

### Rationale for Using Lambda Layers over Docker

While it is technically possible to use Docker containers to bundle Chromium and other
dependencies, there are a few compelling reasons why I chose AWS Lambda Layers over Docker:

1. **Simplicity**: Introducing Docker means adding another layer of complexity. To be
proficient, one needs to understand Docker concepts, write and optimize Dockerfiles,
and potentially grapple with issues related to container orchestration.

2. **Maintenance**: Using Docker implies that there's a need to maintain the Docker
images. As software evolves, the Docker image needs to be updated, possibly re-tested
and then re-deployed. This introduces an ongoing commitment to upkeep.

3. **Faster Deployment**: Deployments using AWS Lambda Layers are considerably faster
compared to deploying Docker containers. This speed becomes especially valuable during
iterative development or frequent deployments.

In summary, while Docker offers great flexibility, its advantages weren't necessary for
this project, and its overhead would've outweighed its benefits.
Thus, AWS Lambda Layers proved to be a more efficient and streamlined choice for our needs.

## Running Locally

If you'd like to test and run `go-screenshoter` on your local machine, you can easily
do so with the following command:

```bash
go run main.go --url=<URL> --path=<screenshot_path>
```

Replace `<URL>` with the website URL you'd like to capture, and `<screenshot_path>`
with the local path where you'd like the screenshot to be saved.

For example:

```bash
go run main.go --url=https://www.google.com --path=lib/google.com.png
```

This provides a quick way to test the functionality before deploying it or to simply
use the tool for local screenshot capturing tasks.

For testing specific parts of the code you should use unit tests and change them to
your needs.

Keep in mind that if you want to test the uploading to R2 you'll need to manually add
the following variables to the `.env` file:

```
CF_ACCOUNT_ID=...
CF_ACCESS_KEY_ID=...
CF_ACCESS_KEY_SECRET=...
CF_BUCKET_NAME=...
```

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
