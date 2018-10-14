# Nibbler Redis

All connection URLs are to not have the protocol in front (i.e. no "http"/"https").

If "Url" is not set on the extension instance, it will pull DB connection info and credentials from REDIS_URL. 
If not found, it will try DATABASE_URL.

If "Password" is not set on the extension instance, it will pull it from REDIS_PASSWORD.  If not found, it will try
DATABASE_PASSWORD.

