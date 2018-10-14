# Nibbler Mongo

A Nibbler Mongo extension.  Manages making an initial connection, then making a client available.

## Configuration

If "Url" is not set on the extension instance, it will pull DB connection info and credentials from MONGO_URL. 
If not found, it will try DATABASE_URL.

