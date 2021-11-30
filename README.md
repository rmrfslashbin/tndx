# tndx
Twitter Indexer

## Summary
tndx is a proof of concept project to learn Go. The operations of the project fetches user details, timeline, and favorites from Twitter then converts to JSON format, compresses, and stores in an AWS S3 bucket.

## AWS Services
### SSM Parameter Store
Configuration data including Twitter API keys and S3 bucket details are stores in AWS SSM the Parameter Store. A local .aws credentials file must be set up to access the data.

### S3
JSON data files are uploaded to an S3 bucket specified in the SSM Param Store. The files are stored in such a way to allow AWS Glue and Athena to "crawl" the data and then query the resulting tables.

## Twitter
A new Twitter project and application were configured for this project. API/Consumer and Access keys/secrets are stored in the AWS SSM Param Store. Further development of this module may include activation of OAuth2 services to query Twitter as other users.

## Subcommands
The Main package of this modules leverages a sysetm of sub-command parsers. This allows the command to be run with a specified sub-command for execution.

## DDB Config
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