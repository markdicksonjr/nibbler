# Nibbler User Auth0

Connects Auth0 and our user model

## Sample app

- should return 404 / "not found" when hitting localhost:<port>/test
- should redirect to login when hitting localhost:<port>/login
- once logged in, it redirects back to localhost:<port>/callback, it will show 404, which is fine
- go to localhost:<port>/test - you should see "authorized"
