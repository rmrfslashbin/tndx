CREATE EXTERNAL TABLE `favorites`(
  `coordinates` string COMMENT 'from deserializer', 
  `created_at` string COMMENT 'from deserializer', 
  `current_user_retweet` string COMMENT 'from deserializer', 
  `entities` struct<hashtags:array<struct<indices:array<int>,text:string>>,media:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string,id:bigint,id_str:string,media_url:string,media_url_https:string,source_status_id:bigint,source_status_id_str:string,type:string,sizes:struct<thumb:struct<w:int,h:int,resize:string>,large:struct<w:int,h:int,resize:string>,medium:struct<w:int,h:int,resize:string>,small:struct<w:int,h:int,resize:string>>,video_info:struct<aspect_ratio:array<int>,duration_millis:int,variants:string>>>,urls:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string>>,user_mentions:array<struct<indices:array<int>,id:bigint,id_str:string,name:string,screen_name:string>>> COMMENT 'from deserializer', 
  `favorite_count` int COMMENT 'from deserializer', 
  `favorited` boolean COMMENT 'from deserializer', 
  `filter_level` string COMMENT 'from deserializer', 
  `id` bigint COMMENT 'from deserializer', 
  `id_str` string COMMENT 'from deserializer', 
  `in_reply_to_screen_name` string COMMENT 'from deserializer', 
  `in_reply_to_status_id` bigint COMMENT 'from deserializer', 
  `in_reply_to_status_id_str` string COMMENT 'from deserializer', 
  `in_reply_to_user_id` bigint COMMENT 'from deserializer', 
  `in_reply_to_user_id_str` string COMMENT 'from deserializer', 
  `lang` string COMMENT 'from deserializer', 
  `possibly_sensitive` boolean COMMENT 'from deserializer', 
  `quote_count` int COMMENT 'from deserializer', 
  `reply_count` int COMMENT 'from deserializer', 
  `retweet_count` int COMMENT 'from deserializer', 
  `retweeted` boolean COMMENT 'from deserializer', 
  `retweeted_status` string COMMENT 'from deserializer', 
  `source` string COMMENT 'from deserializer', 
  `scopes` string COMMENT 'from deserializer', 
  `text` string COMMENT 'from deserializer', 
  `full_text` string COMMENT 'from deserializer', 
  `display_text_range` array<int> COMMENT 'from deserializer', 
  `place` struct<attributes:string,bounding_box:struct<coordinates:array<array<array<double>>>,type:string>,country:string,country_code:string,full_name:string,geometry:string,id:string,name:string,place_type:string,polylines:string,url:string> COMMENT 'from deserializer', 
  `truncated` boolean COMMENT 'from deserializer', 
  `user` struct<contributors_enabled:boolean,created_at:string,default_profile:boolean,default_profile_image:boolean,description:string,email:string,entities:struct<url:struct<hashtags:string,media:string,urls:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string>>,user_mentions:string>,description:struct<hashtags:string,media:string,urls:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string>>,user_mentions:string>>,favourites_count:int,follow_request_sent:boolean,following:boolean,followers_count:int,friends_count:int,geo_enabled:boolean,id:bigint,id_str:string,is_translator:boolean,lang:string,listed_count:int,location:string,name:string,notifications:boolean,profile_background_color:string,profile_background_image_url:string,profile_background_image_url_https:string,profile_background_tile:boolean,profile_banner_url:string,profile_image_url:string,profile_image_url_https:string,profile_link_color:string,profile_sidebar_border_color:string,profile_sidebar_fill_color:string,profile_text_color:string,profile_use_background_image:boolean,protected:boolean,screen_name:string,show_all_inline_media:boolean,status:string,statuses_count:int,time_zone:string,url:string,utc_offset:int,verified:boolean,withheld_in_countries:array<string>,withheld_scope:string> COMMENT 'from deserializer', 
  `withheld_copyright` boolean COMMENT 'from deserializer', 
  `withheld_in_countries` string COMMENT 'from deserializer', 
  `withheld_scope` string COMMENT 'from deserializer', 
  `extended_entities` struct<media:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string,id:bigint,id_str:string,media_url:string,media_url_https:string,source_status_id:bigint,source_status_id_str:string,type:string,sizes:struct<thumb:struct<w:int,h:int,resize:string>,large:struct<w:int,h:int,resize:string>,medium:struct<w:int,h:int,resize:string>,small:struct<w:int,h:int,resize:string>>,video_info:struct<aspect_ratio:array<int>,duration_millis:int,variants:array<struct<content_type:string,bitrate:int,url:string>>>>>> COMMENT 'from deserializer', 
  `extended_tweet` string COMMENT 'from deserializer', 
  `quoted_status_id` bigint COMMENT 'from deserializer', 
  `quoted_status_id_str` string COMMENT 'from deserializer', 
  `quoted_status` struct<coordinates:string,created_at:string,current_user_retweet:string,entities:struct<hashtags:array<struct<indices:array<int>,text:string>>,media:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string,id:bigint,id_str:string,media_url:string,media_url_https:string,source_status_id:bigint,source_status_id_str:string,type:string,sizes:struct<thumb:struct<w:int,h:int,resize:string>,large:struct<w:int,h:int,resize:string>,medium:struct<w:int,h:int,resize:string>,small:struct<w:int,h:int,resize:string>>,video_info:struct<aspect_ratio:array<int>,duration_millis:int,variants:string>>>,urls:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string>>,user_mentions:array<struct<indices:array<int>,id:bigint,id_str:string,name:string,screen_name:string>>>,favorite_count:int,favorited:boolean,filter_level:string,id:bigint,id_str:string,in_reply_to_screen_name:string,in_reply_to_status_id:bigint,in_reply_to_status_id_str:string,in_reply_to_user_id:bigint,in_reply_to_user_id_str:string,lang:string,possibly_sensitive:boolean,quote_count:int,reply_count:int,retweet_count:int,retweeted:boolean,retweeted_status:string,source:string,scopes:string,text:string,full_text:string,display_text_range:array<int>,place:struct<attributes:string,bounding_box:struct<coordinates:array<array<array<double>>>,type:string>,country:string,country_code:string,full_name:string,geometry:string,id:string,name:string,place_type:string,polylines:string,url:string>,truncated:boolean,user:struct<contributors_enabled:boolean,created_at:string,default_profile:boolean,default_profile_image:boolean,description:string,email:string,entities:struct<url:struct<hashtags:string,media:string,urls:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string>>,user_mentions:string>,description:struct<hashtags:string,media:string,urls:array<string>,user_mentions:string>>,favourites_count:int,follow_request_sent:boolean,following:boolean,followers_count:int,friends_count:int,geo_enabled:boolean,id:bigint,id_str:string,is_translator:boolean,lang:string,listed_count:int,location:string,name:string,notifications:boolean,profile_background_color:string,profile_background_image_url:string,profile_background_image_url_https:string,profile_background_tile:boolean,profile_banner_url:string,profile_image_url:string,profile_image_url_https:string,profile_link_color:string,profile_sidebar_border_color:string,profile_sidebar_fill_color:string,profile_text_color:string,profile_use_background_image:boolean,protected:boolean,screen_name:string,show_all_inline_media:boolean,status:string,statuses_count:int,time_zone:string,url:string,utc_offset:int,verified:boolean,withheld_in_countries:array<string>,withheld_scope:string>,withheld_copyright:boolean,withheld_in_countries:string,withheld_scope:string,extended_entities:struct<media:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string,id:bigint,id_str:string,media_url:string,media_url_https:string,source_status_id:bigint,source_status_id_str:string,type:string,sizes:struct<thumb:struct<w:int,h:int,resize:string>,large:struct<w:int,h:int,resize:string>,medium:struct<w:int,h:int,resize:string>,small:struct<w:int,h:int,resize:string>>,video_info:struct<aspect_ratio:array<int>,duration_millis:int,variants:array<struct<content_type:string,bitrate:int,url:string>>>>>>,extended_tweet:string,quoted_status_id:bigint,quoted_status_id_str:string,quoted_status:string> COMMENT 'from deserializer')
PARTITIONED BY ( 
  `partition_0` string)
ROW FORMAT SERDE 
  'org.openx.data.jsonserde.JsonSerDe' 
WITH SERDEPROPERTIES ( 
  'paths'='coordinates,created_at,current_user_retweet,display_text_range,entities,extended_entities,extended_tweet,favorite_count,favorited,filter_level,full_text,id,id_str,in_reply_to_screen_name,in_reply_to_status_id,in_reply_to_status_id_str,in_reply_to_user_id,in_reply_to_user_id_str,lang,place,possibly_sensitive,quote_count,quoted_status,quoted_status_id,quoted_status_id_str,reply_count,retweet_count,retweeted,retweeted_status,scopes,source,text,truncated,user,withheld_copyright,withheld_in_countries,withheld_scope') 
STORED AS INPUTFORMAT 
  'org.apache.hadoop.mapred.TextInputFormat' 
OUTPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
LOCATION
  's3://is-tndx-us-east-2/favorites/'
TBLPROPERTIES (
  'CrawlerSchemaDeserializerVersion'='1.0', 
  'CrawlerSchemaSerializerVersion'='1.0', 
  'UPDATED_BY_CRAWLER'='tndx-favorites', 
  'averageRecordSize'='3864', 
  'classification'='json', 
  'compressionType'='gzip', 
  'objectCount'='198', 
  'recordCount'='198', 
  'sizeKey'='276023', 
  'typeOfData'='file')


  CREATE EXTERNAL TABLE `followers`(
  `contributors_enabled` boolean COMMENT 'from deserializer', 
  `created_at` string COMMENT 'from deserializer', 
  `default_profile` boolean COMMENT 'from deserializer', 
  `default_profile_image` boolean COMMENT 'from deserializer', 
  `description` string COMMENT 'from deserializer', 
  `email` string COMMENT 'from deserializer', 
  `entities` struct<url:struct<hashtags:string,media:string,urls:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string>>,user_mentions:string>,description:struct<hashtags:string,media:string,urls:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string>>,user_mentions:string>> COMMENT 'from deserializer', 
  `favourites_count` int COMMENT 'from deserializer', 
  `follow_request_sent` boolean COMMENT 'from deserializer', 
  `following` boolean COMMENT 'from deserializer', 
  `followers_count` int COMMENT 'from deserializer', 
  `friends_count` int COMMENT 'from deserializer', 
  `geo_enabled` boolean COMMENT 'from deserializer', 
  `id` bigint COMMENT 'from deserializer', 
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
  `status` struct<coordinates:string,created_at:string,current_user_retweet:string,entities:struct<hashtags:array<struct<indices:array<int>,text:string>>,media:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string,id:bigint,id_str:string,media_url:string,media_url_https:string,source_status_id:bigint,source_status_id_str:string,type:string,sizes:struct<thumb:struct<w:int,h:int,resize:string>,large:struct<w:int,h:int,resize:string>,medium:struct<w:int,h:int,resize:string>,small:struct<w:int,h:int,resize:string>>,video_info:struct<aspect_ratio:array<int>,duration_millis:int,variants:string>>>,urls:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string>>,user_mentions:array<struct<indices:array<int>,id:bigint,id_str:string,name:string,screen_name:string>>>,favorite_count:int,favorited:boolean,filter_level:string,id:bigint,id_str:string,in_reply_to_screen_name:string,in_reply_to_status_id:bigint,in_reply_to_status_id_str:string,in_reply_to_user_id:bigint,in_reply_to_user_id_str:string,lang:string,possibly_sensitive:boolean,quote_count:int,reply_count:int,retweet_count:int,retweeted:boolean,retweeted_status:struct<coordinates:string,created_at:string,current_user_retweet:string,entities:struct<hashtags:array<struct<indices:array<int>,text:string>>,media:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string,id:bigint,id_str:string,media_url:string,media_url_https:string,source_status_id:int,source_status_id_str:string,type:string,sizes:struct<thumb:struct<w:int,h:int,resize:string>,large:struct<w:int,h:int,resize:string>,medium:struct<w:int,h:int,resize:string>,small:struct<w:int,h:int,resize:string>>,video_info:struct<aspect_ratio:array<int>,duration_millis:int,variants:string>>>,urls:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string>>,user_mentions:array<struct<indices:array<int>,id:bigint,id_str:string,name:string,screen_name:string>>>,favorite_count:int,favorited:boolean,filter_level:string,id:bigint,id_str:string,in_reply_to_screen_name:string,in_reply_to_status_id:bigint,in_reply_to_status_id_str:string,in_reply_to_user_id:int,in_reply_to_user_id_str:string,lang:string,possibly_sensitive:boolean,quote_count:int,reply_count:int,retweet_count:int,retweeted:boolean,retweeted_status:string,source:string,scopes:string,text:string,full_text:string,display_text_range:array<int>,place:struct<attributes:string,bounding_box:struct<coordinates:array<array<array<double>>>,type:string>,country:string,country_code:string,full_name:string,geometry:string,id:string,name:string,place_type:string,polylines:string,url:string>,truncated:boolean,user:string,withheld_copyright:boolean,withheld_in_countries:string,withheld_scope:string,extended_entities:struct<media:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string,id:bigint,id_str:string,media_url:string,media_url_https:string,source_status_id:int,source_status_id_str:string,type:string,sizes:struct<thumb:struct<w:int,h:int,resize:string>,large:struct<w:int,h:int,resize:string>,medium:struct<w:int,h:int,resize:string>,small:struct<w:int,h:int,resize:string>>,video_info:struct<aspect_ratio:array<int>,duration_millis:int,variants:array<struct<content_type:string,bitrate:int,url:string>>>>>>,extended_tweet:string,quoted_status_id:bigint,quoted_status_id_str:string,quoted_status:string>,source:string,scopes:string,text:string,full_text:string,display_text_range:array<int>,place:struct<attributes:string,bounding_box:struct<coordinates:array<array<array<double>>>,type:string>,country:string,country_code:string,full_name:string,geometry:string,id:string,name:string,place_type:string,polylines:string,url:string>,truncated:boolean,user:string,withheld_copyright:boolean,withheld_in_countries:string,withheld_scope:string,extended_entities:struct<media:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string,id:bigint,id_str:string,media_url:string,media_url_https:string,source_status_id:bigint,source_status_id_str:string,type:string,sizes:struct<thumb:struct<w:int,h:int,resize:string>,large:struct<w:int,h:int,resize:string>,medium:struct<w:int,h:int,resize:string>,small:struct<w:int,h:int,resize:string>>,video_info:struct<aspect_ratio:array<int>,duration_millis:int,variants:array<struct<content_type:string,bitrate:int,url:string>>>>>>,extended_tweet:string,quoted_status_id:bigint,quoted_status_id_str:string,quoted_status:string> COMMENT 'from deserializer', 
  `statuses_count` int COMMENT 'from deserializer', 
  `time_zone` string COMMENT 'from deserializer', 
  `url` string COMMENT 'from deserializer', 
  `utc_offset` int COMMENT 'from deserializer', 
  `verified` boolean COMMENT 'from deserializer', 
  `withheld_in_countries` array<string> COMMENT 'from deserializer', 
  `withheld_scope` string COMMENT 'from deserializer')
PARTITIONED BY ( 
  `partition_0` string)
ROW FORMAT SERDE 
  'org.openx.data.jsonserde.JsonSerDe' 
WITH SERDEPROPERTIES ( 
  'paths'='contributors_enabled,created_at,default_profile,default_profile_image,description,email,entities,favourites_count,follow_request_sent,followers_count,following,friends_count,geo_enabled,id,id_str,is_translator,lang,listed_count,location,name,notifications,profile_background_color,profile_background_image_url,profile_background_image_url_https,profile_background_tile,profile_banner_url,profile_image_url,profile_image_url_https,profile_link_color,profile_sidebar_border_color,profile_sidebar_fill_color,profile_text_color,profile_use_background_image,protected,screen_name,show_all_inline_media,status,statuses_count,time_zone,url,utc_offset,verified,withheld_in_countries,withheld_scope') 
STORED AS INPUTFORMAT 
  'org.apache.hadoop.mapred.TextInputFormat' 
OUTPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
LOCATION
  's3://is-tndx-us-east-2/followers/'
TBLPROPERTIES (
  'CrawlerSchemaDeserializerVersion'='1.0', 
  'CrawlerSchemaSerializerVersion'='1.0', 
  'UPDATED_BY_CRAWLER'='tndx-followers', 
  'averageRecordSize'='3385', 
  'classification'='json', 
  'compressionType'='gzip', 
  'objectCount'='117', 
  'recordCount'='117', 
  'sizeKey'='144983', 
  'typeOfData'='file')




  CREATE EXTERNAL TABLE `friends`(
  `contributors_enabled` boolean COMMENT 'from deserializer', 
  `created_at` string COMMENT 'from deserializer', 
  `default_profile` boolean COMMENT 'from deserializer', 
  `default_profile_image` boolean COMMENT 'from deserializer', 
  `description` string COMMENT 'from deserializer', 
  `email` string COMMENT 'from deserializer', 
  `entities` struct<url:struct<hashtags:string,media:string,urls:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string>>,user_mentions:string>,description:struct<hashtags:string,media:string,urls:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string>>,user_mentions:string>> COMMENT 'from deserializer', 
  `favourites_count` int COMMENT 'from deserializer', 
  `follow_request_sent` boolean COMMENT 'from deserializer', 
  `following` boolean COMMENT 'from deserializer', 
  `followers_count` int COMMENT 'from deserializer', 
  `friends_count` int COMMENT 'from deserializer', 
  `geo_enabled` boolean COMMENT 'from deserializer', 
  `id` bigint COMMENT 'from deserializer', 
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
  `status` struct<coordinates:struct<coordinates:array<double>,type:string>,created_at:string,current_user_retweet:string,entities:struct<hashtags:array<struct<indices:array<int>,text:string>>,media:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string,id:bigint,id_str:string,media_url:string,media_url_https:string,source_status_id:bigint,source_status_id_str:string,type:string,sizes:struct<thumb:struct<w:int,h:int,resize:string>,large:struct<w:int,h:int,resize:string>,medium:struct<w:int,h:int,resize:string>,small:struct<w:int,h:int,resize:string>>,video_info:struct<aspect_ratio:array<int>,duration_millis:int,variants:string>>>,urls:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string>>,user_mentions:array<struct<indices:array<int>,id:bigint,id_str:string,name:string,screen_name:string>>>,favorite_count:int,favorited:boolean,filter_level:string,id:bigint,id_str:string,in_reply_to_screen_name:string,in_reply_to_status_id:bigint,in_reply_to_status_id_str:string,in_reply_to_user_id:bigint,in_reply_to_user_id_str:string,lang:string,possibly_sensitive:boolean,quote_count:int,reply_count:int,retweet_count:int,retweeted:boolean,retweeted_status:struct<coordinates:struct<coordinates:array<double>,type:string>,created_at:string,current_user_retweet:string,entities:struct<hashtags:array<struct<indices:array<int>,text:string>>,media:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string,id:bigint,id_str:string,media_url:string,media_url_https:string,source_status_id:bigint,source_status_id_str:string,type:string,sizes:struct<thumb:struct<w:int,h:int,resize:string>,large:struct<w:int,h:int,resize:string>,medium:struct<w:int,h:int,resize:string>,small:struct<w:int,h:int,resize:string>>,video_info:struct<aspect_ratio:array<int>,duration_millis:int,variants:string>>>,urls:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string>>,user_mentions:array<struct<indices:array<int>,id:bigint,id_str:string,name:string,screen_name:string>>>,favorite_count:int,favorited:boolean,filter_level:string,id:bigint,id_str:string,in_reply_to_screen_name:string,in_reply_to_status_id:bigint,in_reply_to_status_id_str:string,in_reply_to_user_id:bigint,in_reply_to_user_id_str:string,lang:string,possibly_sensitive:boolean,quote_count:int,reply_count:int,retweet_count:int,retweeted:boolean,retweeted_status:string,source:string,scopes:string,text:string,full_text:string,display_text_range:array<int>,place:struct<attributes:string,bounding_box:struct<coordinates:array<array<array<double>>>,type:string>,country:string,country_code:string,full_name:string,geometry:string,id:string,name:string,place_type:string,polylines:string,url:string>,truncated:boolean,user:string,withheld_copyright:boolean,withheld_in_countries:string,withheld_scope:string,extended_entities:struct<media:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string,id:bigint,id_str:string,media_url:string,media_url_https:string,source_status_id:bigint,source_status_id_str:string,type:string,sizes:struct<thumb:struct<w:int,h:int,resize:string>,large:struct<w:int,h:int,resize:string>,medium:struct<w:int,h:int,resize:string>,small:struct<w:int,h:int,resize:string>>,video_info:struct<aspect_ratio:array<int>,duration_millis:int,variants:array<struct<content_type:string,bitrate:int,url:string>>>>>>,extended_tweet:string,quoted_status_id:bigint,quoted_status_id_str:string,quoted_status:string>,source:string,scopes:string,text:string,full_text:string,display_text_range:array<int>,place:struct<attributes:string,bounding_box:struct<coordinates:array<array<array<double>>>,type:string>,country:string,country_code:string,full_name:string,geometry:string,id:string,name:string,place_type:string,polylines:string,url:string>,truncated:boolean,user:string,withheld_copyright:boolean,withheld_in_countries:string,withheld_scope:string,extended_entities:struct<media:array<struct<indices:array<int>,display_url:string,expanded_url:string,url:string,id:bigint,id_str:string,media_url:string,media_url_https:string,source_status_id:bigint,source_status_id_str:string,type:string,sizes:struct<thumb:struct<w:int,h:int,resize:string>,large:struct<w:int,h:int,resize:string>,medium:struct<w:int,h:int,resize:string>,small:struct<w:int,h:int,resize:string>>,video_info:struct<aspect_ratio:array<int>,duration_millis:int,variants:array<struct<content_type:string,bitrate:int,url:string>>>>>>,extended_tweet:string,quoted_status_id:bigint,quoted_status_id_str:string,quoted_status:string> COMMENT 'from deserializer', 
  `statuses_count` int COMMENT 'from deserializer', 
  `time_zone` string COMMENT 'from deserializer', 
  `url` string COMMENT 'from deserializer', 
  `utc_offset` int COMMENT 'from deserializer', 
  `verified` boolean COMMENT 'from deserializer', 
  `withheld_in_countries` array<string> COMMENT 'from deserializer', 
  `withheld_scope` string COMMENT 'from deserializer')
PARTITIONED BY ( 
  `partition_0` string)
ROW FORMAT SERDE 
  'org.openx.data.jsonserde.JsonSerDe' 
WITH SERDEPROPERTIES ( 
  'paths'='contributors_enabled,created_at,default_profile,default_profile_image,description,email,entities,favourites_count,follow_request_sent,followers_count,following,friends_count,geo_enabled,id,id_str,is_translator,lang,listed_count,location,name,notifications,profile_background_color,profile_background_image_url,profile_background_image_url_https,profile_background_tile,profile_banner_url,profile_image_url,profile_image_url_https,profile_link_color,profile_sidebar_border_color,profile_sidebar_fill_color,profile_text_color,profile_use_background_image,protected,screen_name,show_all_inline_media,status,statuses_count,time_zone,url,utc_offset,verified,withheld_in_countries,withheld_scope') 
STORED AS INPUTFORMAT 
  'org.apache.hadoop.mapred.TextInputFormat' 
OUTPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
LOCATION
  's3://is-tndx-us-east-2/friends/'
TBLPROPERTIES (
  'CrawlerSchemaDeserializerVersion'='1.0', 
  'CrawlerSchemaSerializerVersion'='1.0', 
  'UPDATED_BY_CRAWLER'='tndx-friends', 
  'averageRecordSize'='3346', 
  'classification'='json', 
  'compressionType'='gzip', 
  'objectCount'='1031', 
  'recordCount'='1031', 
  'sizeKey'='1342619', 
  'typeOfData'='file')




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