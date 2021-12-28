  CREATE EXTERNAL TABLE `users`(
  `contributors_enabled` boolean COMMENT 'from deserializer', 
  `created_at` string COMMENT 'from deserializer', 
  `default_profile` boolean COMMENT 'from deserializer', 
  `default_profile_image` boolean COMMENT 'from deserializer', 
  `description` string COMMENT 'from deserializer', 
  `email` string COMMENT 'from deserializer', 
  `entities` struct<url:struct<hashtags:string,media:string,urls:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string>>,user_mentions:string>,description:struct<hashtags:string,media:string,urls:array<string>,user_mentions:string>> COMMENT 'from deserializer', 
  `favourites_count` int COMMENT 'from deserializer', 
  `follow_request_sent` boolean COMMENT 'from deserializer', 
  `following` boolean COMMENT 'from deserializer', 
  `followers_count` int COMMENT 'from deserializer', 
  `friends_count` int COMMENT 'from deserializer', 
  `geo_enabled` boolean COMMENT 'from deserializer', 
  `id` int COMMENT 'from deserializer', 
  `id_str` string COMMENT 'from deserializer', 
  `is_translator` boolean COMMENT 'from deserializer', 
  `lang` string COMMENT 'from deserializer', 
  `listed_count` int COMMENT 'from deserializer', 
  `location` string COMMENT 'from deserializer', 
  `name` string COMMENT 'from deserializer', 
  `notifications` boolean COMMENT 'from deserializer', 
  `profile_background_color` string COMMENT 'from deserializer', 
  `profile_background_image_url` string COMMENT 'from deserializer', 
  `profile_background_image_url_https` string COMMENT 'from deserializer', 
  `profile_background_tile` boolean COMMENT 'from deserializer', 
  `profile_banner_url` string COMMENT 'from deserializer', 
  `profile_image_url` string COMMENT 'from deserializer', 
  `profile_image_url_https` string COMMENT 'from deserializer', 
  `profile_link_color` string COMMENT 'from deserializer', 
  `profile_sidebar_border_color` string COMMENT 'from deserializer', 
  `profile_sidebar_fill_color` string COMMENT 'from deserializer', 
  `profile_text_color` string COMMENT 'from deserializer', 
  `profile_use_background_image` boolean COMMENT 'from deserializer', 
  `protected` boolean COMMENT 'from deserializer', 
  `screen_name` string COMMENT 'from deserializer', 
  `show_all_inline_media` boolean COMMENT 'from deserializer', 
  `status` struct<coordinates:string,created_at:string,current_user_retweet:string,entities:struct<hashtags:array<struct<indices:array<int>,text:string>>,media:string,urls:array<string>,user_mentions:array<struct<indices:array<int>,id:int,id_str:string,name:string,screen_name:string>>>,favorite_count:int,favorited:boolean,filter_level:string,id:bigint,id_str:string,in_reply_to_screen_name:string,in_reply_to_status_id:int,in_reply_to_status_id_str:string,in_reply_to_user_id:int,in_reply_to_user_id_str:string,lang:string,possibly_sensitive:boolean,quote_count:int,reply_count:int,retweet_count:int,retweeted:boolean,retweeted_status:struct<coordinates:string,created_at:string,current_user_retweet:string,entities:struct<hashtags:array<struct<indices:array<int>,text:string>>,media:string,urls:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string>>,user_mentions:array<string>>,favorite_count:int,favorited:boolean,filter_level:string,id:bigint,id_str:string,in_reply_to_screen_name:string,in_reply_to_status_id:int,in_reply_to_status_id_str:string,in_reply_to_user_id:int,in_reply_to_user_id_str:string,lang:string,possibly_sensitive:boolean,quote_count:int,reply_count:int,retweet_count:int,retweeted:boolean,retweeted_status:string,source:string,scopes:string,text:string,full_text:string,display_text_range:array<int>,place:string,truncated:boolean,user:string,withheld_copyright:boolean,withheld_in_countries:string,withheld_scope:string,extended_entities:string,extended_tweet:string,quoted_status_id:int,quoted_status_id_str:string,quoted_status:string>,source:string,scopes:string,text:string,full_text:string,display_text_range:array<int>,place:string,truncated:boolean,user:string,withheld_copyright:boolean,withheld_in_countries:string,withheld_scope:string,extended_entities:string,extended_tweet:string,quoted_status_id:int,quoted_status_id_str:string,quoted_status:string> COMMENT 'from deserializer', 
  `statuses_count` int COMMENT 'from deserializer', 
  `time_zone` string COMMENT 'from deserializer', 
  `url` string COMMENT 'from deserializer', 
  `utc_offset` int COMMENT 'from deserializer', 
  `verified` boolean COMMENT 'from deserializer', 
  `withheld_in_countries` array<string> COMMENT 'from deserializer', 
  `withheld_scope` string COMMENT 'from deserializer')
ROW FORMAT SERDE 
  'org.openx.data.jsonserde.JsonSerDe' 
WITH SERDEPROPERTIES ( 
  'paths'='contributors_enabled,created_at,default_profile,default_profile_image,description,email,entities,favourites_count,follow_request_sent,followers_count,following,friends_count,geo_enabled,id,id_str,is_translator,lang,listed_count,location,name,notifications,profile_background_color,profile_background_image_url,profile_background_image_url_https,profile_background_tile,profile_banner_url,profile_image_url,profile_image_url_https,profile_link_color,profile_sidebar_border_color,profile_sidebar_fill_color,profile_text_color,profile_use_background_image,protected,screen_name,show_all_inline_media,status,statuses_count,time_zone,url,utc_offset,verified,withheld_in_countries,withheld_scope') 
STORED AS INPUTFORMAT 
  'org.apache.hadoop.mapred.TextInputFormat' 
OUTPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
LOCATION
  's3://is-tndx-us-east-2/users/'
TBLPROPERTIES (
  'CrawlerSchemaDeserializerVersion'='1.0', 
  'CrawlerSchemaSerializerVersion'='1.0', 
  'UPDATED_BY_CRAWLER'='tndx-users', 
  'averageRecordSize'='4374', 
  'classification'='json', 
  'compressionType'='gzip', 
  'objectCount'='1', 
  'recordCount'='1', 
  'sizeKey'='1406', 
  'typeOfData'='file')