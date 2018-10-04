#nibbler-azure-blob

A simple Azure Blob Storage extension for Nibbler

## Usage

Some Nibbler config values are required:

- azure.blob.account.name (AZURE_BLOB_ACCOUNT_NAME env var)
- azure.blob.account.key (AZURE_BLOB_ACCOUNT_KEY env var)

Upon init, the extension will have an initialized the credentials for Azure.  To use
the extension, call the GetContainerURL method and perform operations on the result.

This extension obviously requires an Azure Storage Account.  Once an account is added,
go to "Access keys" on the storage account.  Copy the account name and one of the keys and
set the appropriate environment variables listed above