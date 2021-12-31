# tndx
Twitter Indexer & Archiver

## Summary
```tndx``` fetches user details, timeline, followers, friends, and favorites from Twitter, then stores in an AWS S3 bucket via a [Kinesis Delivery Stream](https://docs.aws.amazon.com/firehose/latest/dev/what-is-this-service.html) for query via [Athena](https://aws.amazon.com/athena/)/[Trino](https://trino.io). ```tndx``` also extract "entities" (media) URLs from tweets, fetches and processes via [AWS Rekognition](https://aws.amazon.com/rekognition/), storing the media in S3 and resulting media meta data in DynamoDB.

## AWS Services
An AWS account and local [.aws credentials file](https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/setup-credentials.html) must be set up to run ```tndx```.

### SSM Parameter Store
Configuration data including Twitter API keys and S3 bucket details are stored in the [AWS SSM Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html). 

### S3
ORC data files are stored in an [S3](https://aws.amazon.com/s3/) bucket specified in the SSM Param Store. The files are stored in such a way to allow [AWS Glue](https://aws.amazon.com/glue/) and Athena to "crawl" the data and then query the resulting tables.

## Twitter
A [Twitter project and application](https://developer.twitter.com/) must be configured for this project. API/Consumer keys are stored in the AWS SSM Param Store. ```tndx``` uses Twitter's [OAuth 2.0](https://developer.twitter.com/en/docs/authentication/oauth-2-0) services.

