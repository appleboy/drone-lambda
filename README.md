# drone-lambda

[![GoDoc](https://godoc.org/github.com/appleboy/drone-lambda?status.svg)](https://godoc.org/github.com/appleboy/drone-lambda)
[![Build Status](https://cloud.drone.io/api/badges/appleboy/drone-lambda/status.svg)](https://cloud.drone.io/appleboy/drone-lambda)
[![codecov](https://codecov.io/gh/appleboy/drone-lambda/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-lambda)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-lambda)](https://goreportcard.com/report/github.com/appleboy/drone-lambda)
[![Docker Pulls](https://img.shields.io/docker/pulls/appleboy/drone-lambda.svg)](https://hub.docker.com/r/appleboy/drone-lambda/)
[![](https://images.microbadger.com/badges/image/appleboy/drone-lambda.svg)](https://microbadger.com/images/appleboy/drone-lambda "Get your own image badge on microbadger.com")
[![Release](https://github-release-version.herokuapp.com/github/appleboy/drone-lambda/release.svg?style=flat)](https://github.com/appleboy/drone-lambda/releases/latest)
[![Build status](https://ci.appveyor.com/api/projects/status/cuioqombam9yufdy?svg=true)](https://ci.appveyor.com/project/appleboy/drone-lambda)



Deploying Lambda code with drone CI to an existing function

## Build or Download a binary

The pre-compiled binaries can be downloaded from [release page](https://github.com/appleboy/drone-lambda/releases). Support the following OS type.

* Windows amd64/386
* Linux amd64/386
* Darwin amd64/386

With `Go` installed

```
$ go get -u -v github.com/appleboy/drone-lambda
``` 

or build the binary with the following command:

```
$ make build
```

## Docker

Build the docker image with the following commands:

```
$ make docker
```

Please note incorrectly building the image for the correct x64 linux and with
CGO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-lambda' not found or does not exist..
```

## Usage

There are three ways to send notification.

- [drone-lambda](#drone-lambda)
  - [Build or Download a binary](#build-or-download-a-binary)
  - [Docker](#docker)
  - [Usage](#usage)
    - [Usage from binary](#usage-from-binary)
    - [Usage from docker](#usage-from-docker)
    - [Usage from drone ci](#usage-from-drone-ci)
  - [AWS Policy](#aws-policy)

<a name="usage-from-binary"></a>
### Usage from binary

Update lambda function from zip file.

```sh
$ drone-lambda --region ap-southeast-1 \
  --access-key xxxx \
  --secret-key xxxx \
  --function-name upload-s3 \
  --zip-file deployment.zip
```

Update lambda function from s3 object.

```sh
$ drone-lambda --region ap-southeast-1 \
  --access-key xxxx \
  --secret-key xxxx \
  --function-name upload-s3 \
  --s3-bucket some-bucket \
  --s3-key lambda-dir/lambda-project-${DRONE_BUILD_NUMBER}.zip
```

<a name="usage-from-docker"></a>
### Usage from docker

Update lambda function from zip file.

```bash
docker run --rm \
  -e AWS_ACCESS_KEY_ID=xxxxxxx \
  -e AWS_SECRET_ACCESS_KEY=xxxxxxx \
  -e FUNCTION_NAME=upload-s3 \
  -e ZIP_FILE=deployment.zip \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-lambda
```

Update lambda function from s3 object.

```bash
docker run --rm \
  -e AWS_ACCESS_KEY_ID=xxxxxxx \
  -e AWS_SECRET_ACCESS_KEY=xxxxxxx \
  -e FUNCTION_NAME=upload-s3 \
  -e S3_BUCKET=some-bucket \
  -e S3_KEY=lambda-project.zip \
  appleboy/drone-lambda
```

<a name="usage-from-drone-ci"></a>
### Usage from drone ci

Update lambda function, execute from the working directory:

```bash
docker run --rm \
  -e AWS_ACCESS_KEY_ID=xxxxxxx \
  -e AWS_SECRET_ACCESS_KEY=xxxxxxx \
  -e FUNCTION_NAME=upload-s3 \
  -e ZIP_FILE=deployment.zip \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-lambda
```

## AWS Policy

Add the following AWS policy if you want to integrate with CI/CD tools like Jenkins, GitLab Ci or Drone.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "iam:ListRoles",
        "lambda:UpdateFunctionCode",
        "lambda:CreateFunction"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    }
  ]
}
```
