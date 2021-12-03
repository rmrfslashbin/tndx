# tndx
Twitter Indexer

## Summary
```tndx``` fetches user details, timeline, followers, friends, and favorites from Twitter, then converts to JSON format, compresses, and stores in an AWS S3 bucket in an [Athena](https://aws.amazon.com/athena/)/[Trino](https://trino.io) compatible manner.

## AWS Services
An AWS account and local [.aws credentials file](https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/setup-credentials.html) must be set up to run ```tndx```.

### SSM Parameter Store
Configuration data including Twitter API keys and S3 bucket details are stored in the [AWS SSM Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html). 

### S3
JSON data files are uploaded to an [S3](https://aws.amazon.com/s3/) bucket specified in the SSM Param Store. The files are stored in such a way to allow [AWS Glue](https://aws.amazon.com/glue/) and Athena to "crawl" the data and then query the resulting tables.

## Twitter
A [Twitter project and application](https://developer.twitter.com/) must be configured for this project. API/Consumer keys are stored in the AWS SSM Param Store. ```tndx``` uses Twitter's [OAuth 2.0](https://developer.twitter.com/en/docs/authentication/oauth-2-0) services.

## Usage
```
Usage:
   [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  entities    fetch entities
  favorites   fetch the user's favorites
  followers   fetch the user's followers
  friends     fetch the user's friends
  help        Help about any command
  timeline    fetch the user's timeline
  user        lookup user by userid or screenname

Flags:
      --database string        [sqlite|ddb]
      --dotenv string          dotenv path (default "./.env")
  -h, --help                   help for this command
      --localrootpath string   local root path (default "./data")
      --loglevel string        [error|warn|info|debug|trace] (default "info")
      --screenname string      screen name
      --storage string         [local|s3]
      --userid int             user id
  -v, --version                version for this command

Use " [command] --help" for more information about a command.
```

### Example
```
% tndx --database ddb --storage s3 favorites --userid 16020064
INFO[0002] finished getting timeline action="RunTimelineCmd::Done!" count=26 lowerID=1465833642157 upperID=1460899680291 userid=16020064
%
```

## AWS Configs

### DDB Config
```
{
    "Table": {
        "AttributeDefinitions": [
            {
                "AttributeName": "domain",
                "AttributeType": "S"
            },
            {
                "AttributeName": "userid",
                "AttributeType": "N"
            }
        ],
        "TableName": "tndx",
        "KeySchema": [
            {
                "AttributeName": "userid",
                "KeyType": "HASH"
            },
            {
                "AttributeName": "domain",
                "KeyType": "RANGE"
            }
        ],
        "TableStatus": "ACTIVE",
        "CreationDateTime": "2021-11-29T22:28:59.008000-05:00",
        "ProvisionedThroughput": {
            "NumberOfDecreasesToday": 0,
            "ReadCapacityUnits": 0,
            "WriteCapacityUnits": 0
        },
        "TableSizeBytes": 123,
        "ItemCount": 2,
        "TableArn": "arn:aws:dynamodb:us-east-1:150319663043:table/tndx",
        "TableId": "19bd48e5-ce49-443c-9b62-ca79e9b3590e",
        "BillingModeSummary": {
            "BillingMode": "PAY_PER_REQUEST",
            "LastUpdateToPayPerRequestDateTime": "2021-11-29T22:28:59.008000-05:00"
        }
    }
}
```

### Athena Table
While ```tndx``` saves the export data in an Glue/Athena mannor, the below resources helped identify proper orginization and output formatting.
* https://aws.amazon.com/premiumsupport/knowledge-center/error-json-athena/
* https://ai-services.go-aws.com/30_social-media-analytics/50_create_athena_tables.html

### Glue
Once data files are saved in S3, use an AWS Glue Crawler to populate Hive data (this is a step towards using Trino to query the saved data). The below items are a general reference for setting up a crawler.
* Go to https://console.aws.amazon.com/glue/home?region=us-east-1#catalog:tab=crawlers
* Set up a new crawler...
* Crawler source type: Data stores
* Repeat crawls of S3 data stores: Crawl all folders
* Include Path: s3://my-tndx-bucket/timeline
* Grouping behavior for S3 data (optional): (check) Create a single schema for each S3 path

Run the crawler and query with Athena or Trino.
