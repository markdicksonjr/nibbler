# Nibbler Session

Provides a session for the server to use.

An SQL connector is available out of the box, which can be used in memory mode by
not providing a DB reference.

The default MaxAge is 30 days (86400 * 30) in all cases (which is pretty long).
