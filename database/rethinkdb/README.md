# Nibbler Redis

All connection URLs have the protocol in front removed during processing (e.g. no "http"/"https").

If "Url" is not set on the extension instance, it will pull DB connection info and credentials from REDIS_URL. 
If not found, it will try REDISCLOUD_URL, then DATABASE_URL.

