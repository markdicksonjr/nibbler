# Nibbler SQL

Requires models that are to be migrated to be provided before init.

Assumes models use GORM, currently.  Future enhancements will guarantee
that GORM models will work for non-GORM-supported SQL DBs (if that's a thing).

Will pull DB connection info and credentials from SQL_URL.  If not found,
it will try DATABASE_URL, then DB_URL.

If any part of the needed credentials are not obtained at that stage, it will use:

`DB_DIALECT`
`DB_USER`
`DB_PASSWORD`
`DB_DBNAME`

At present, the dialects available are limited to `postgres`, `sqlite3`.  If none is detected,
`sqlite3` will be used (in memory).  For example, you could use sqlite3 to connect to a file with `SQL_URL=sqlite3:///tmp/test.db`
or something like `postgres://test:test@localhost:5432/test` to connect via postgres.