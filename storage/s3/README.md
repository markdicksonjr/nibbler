#nibbler-s3

Some Nibbler config values are required:

- s3.accesskey (S3_ACCESSKEY env var)
- s3.secret (S3_SECRET env var)
- s3.endpoint (S3_ENDPOINT env var)
- s3.region (S3_REGION env var)

This was tested using Digital Ocean Spaces, but it should work just as well with AWS S3

## Usage

Upon init, the extension will have an initialized S3 client as the "S3" property