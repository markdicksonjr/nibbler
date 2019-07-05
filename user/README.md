# Nibbler User

Provides a basic user model, and some means to persist and query it.

## Utilities

- GetSafeUser: returns a version of the user with sensitive properties wiped.
- EnforceLoggedIn
- EnforceEmailValidated

## Context and Protected Context

- Some "room" is available in the default user model for app-specific data that
is attached to users.  These are the context properties.  The difference between Context
and Protected Context is whether or not API requests should expose the data or not.  The
Protected Context is a place where the back-end can attach additional info without it ever
being seen by the user, essentially.

