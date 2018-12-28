# Nibbler Session

Provides a session for the server to use.

Users can provide their own connector, or a Cookie connector will be used
if one is not provided.  It is recommended that users use their own connector
for production apps.  A sample has been provided in ./sample.connector.

The default MaxAge is 30 days (86400 * 30) in all cases (which is pretty long).
