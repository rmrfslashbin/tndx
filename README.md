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

## Athena Table
https://aws.amazon.com/premiumsupport/knowledge-center/error-json-athena/
https://ai-services.go-aws.com/30_social-media-analytics/50_create_athena_tables.html

```
CREATE EXTERNAL TABLE tndx.timelines (
    coordinates STRUCT<
        type: STRING,
        coordinates: ARRAY<
            DOUBLE
        >
    >,
    retweeted BOOLEAN,
    source STRING,
    entities STRUCT<
        hashtags: ARRAY<
            STRUCT<
                text: STRING,
                indices: ARRAY<
                    BIGINT
                >
            >
        >,
        urls: ARRAY<
            STRUCT<
                url: STRING,
                expanded_url: STRING,
                display_url: STRING,
                indices: ARRAY<
                    BIGINT
                >
            >
        >
    >,
    reply_count BIGINT,
    favorite_count BIGINT,
    geo STRUCT<
        type: STRING,
        coordinates: ARRAY<
            DOUBLE
        >
    >,
    id_str STRING,
    truncated BOOLEAN,
    text STRING,
    retweet_count BIGINT,
    id BIGINT,
    possibly_sensitive BOOLEAN,
    filter_level STRING,
    created_at STRING,
    place STRUCT<
        id: STRING,
        url: STRING,
        place_type: STRING,
        name: STRING,
        full_name: STRING,
        country_code: STRING,
        country: STRING,
        bounding_box: STRUCT<
            type: STRING,
            coordinates: ARRAY<
                ARRAY<
                    ARRAY<
                        FLOAT
                    >
                >
            >
        >
    >,
    favorited BOOLEAN,
    lang STRING,
    in_reply_to_screen_name STRING,
    is_quote_status BOOLEAN,
    in_reply_to_user_id_str STRING,
    user STRUCT<
        id: BIGINT,
        id_str: STRING,
        name: STRING,
        screen_name: STRING,
        location: STRING,
        url: STRING,
        description: STRING,
        translator_type: STRING,
        protected: BOOLEAN,
        verified: BOOLEAN,
        followers_count: BIGINT,
        friends_count: BIGINT,
        listed_count: BIGINT,
        favourites_count: BIGINT,
        statuses_count: BIGINT,
        created_at: STRING,
        utc_offset: BIGINT,
        time_zone: STRING,
        geo_enabled: BOOLEAN,
        lang: STRING,
        contributors_enabled: BOOLEAN,
        is_translator: BOOLEAN,
        profile_background_color: STRING,
        profile_background_image_url: STRING,
        profile_background_image_url_https: STRING,
        profile_background_tile: BOOLEAN,
        profile_link_color: STRING,
        profile_sidebar_border_color: STRING,
        profile_sidebar_fill_color: STRING,
        profile_text_color: STRING,
        profile_use_background_image: BOOLEAN,
        profile_image_url: STRING,
        profile_image_url_https: STRING,
        profile_banner_url: STRING,
        default_profile: BOOLEAN,
        default_profile_image: BOOLEAN
    >,
    quote_count BIGINT
) ROW FORMAT SERDE 'org.openx.data.jsonserde.JsonSerDe'
LOCATION 's3://you/s3/bucket/path/';
```