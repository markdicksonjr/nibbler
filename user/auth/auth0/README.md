# Nibbler User Auth0

Connects Auth0 and our user model

## Sample app

- first, change sample.application.go and set the user email to match the account in Auth0 (this module
requires the user to be in the local DB)
- should return 404 / "not found" when hitting localhost:<port>/test
- should redirect to login when hitting localhost:<port>/login
- once logged in, it redirects back to localhost:<port>/callback, it will respond with a 200 status code
- go to localhost:<port>/test - you should see "authorized"
