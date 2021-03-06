# Nibbler

[![CircleCI](https://circleci.com/gh/markdicksonjr/nibbler.svg?style=svg)](https://circleci.com/gh/markdicksonjr/nibbler)

An extension-oriented framework designed to take a lot of the boilerplate out of making a top-notch Go web server or 
service worker.  Requires Go v1.12+.  Go v1.13+ recommended. 

A [sample app](https://github.com/markdicksonjr/nibbler-sample) is available that uses Nibbler.  Hopefully it highlights
the ease-of-use of the framework.

## Running the included sample apps

Sample apps have been provided to show how Nibbler and some extensions are used.  They'll be helpful as I fill in many 
documentation gaps.  You can find them in their respective extension repositories.

Build/run the sample app you're interested in using (from the current directory, the one in this repository is in
./sample):

`go run main.go`

## Extensions

Extensions are the backbone of Nibbler.  Extensions must implement a very simple interface (nibbler.Extension).  A base 
class extension (NoOpExtension) is available for cases where only very few Extension methods (or none?) are required for 
the extension you're building.

## Included Extension Categories

Nibbler also provides some extension implementations that perform common tasks for web services.

Extensions provided with Nibbler are organized by functional category, with the main Nibbler structures at the root
level.  These are the module categories below the root level:

- Session - session storage and retrieval.  See ./session/README.md
- User - the Nibbler user model, and various integrations that can operate with it.  These will tend to be auth integrations.
See ./user/README.md
- Local auth - local auth implementation, expects a user extension and session extension implementation to be provided to it.
See ./user/auth/local/README.md

## External Extensions

There are various features provided by external modules:

- Auth - authentication/authorization modules that do not integrate with Nibbler's user model (source of truth is not Nibbler).
- Database storage - connect to databases and expose mechanisms to create, query, etc.
- Mail - outbound email/sms/etc
- Storage - block/blob storage

Here are some repositories containing such extensions (more to come):

https://github.com/markdicksonjr/nibbler-auth0

https://github.com/markdicksonjr/nibbler-elasticsearch

https://github.com/markdicksonjr/nibbler-mail-outbound

https://github.com/markdicksonjr/nibbler-oauth2

https://github.com/markdicksonjr/nibbler-socket

https://github.com/markdicksonjr/nibbler-sql

https://github.com/markdicksonjr/nibbler-storage

## Configuration

The app configuration can be created with this method:

```go
config, err := nibbler.LoadConfiguration()
```

If nil is provided to the core.LoadConfiguration method, it will use environment variables for
configuration of your app, with anything in ./config.json overriding it.  This feature can be overridden 
by doing something like this in your app:

```go
envSources := []source.Source{
    file.NewSource(file.WithPath("./config.json")),
    env.NewSource(),
}

config, err := nibbler.LoadConfiguration(&envSources)
```

In this case (providing nil results in this behavior, but it's shown here explicitly to show how the mechanism works), 
it will first apply environment variables, then apply properties from the json file (reverse order of definition).  The 
environment variables used are upper-case but with dot-notation where the dots are replaced by underscores (e.g. 
nibbler.port is NIBBLER_PORT).

The following properties are available by default, but custom properties can be obtained from your extension from the 
Application passed to it in the Init() method.  For example:

```
app.Configuration.Raw.Get("some", "property")
```

The core properties (all optional) are:

- NIBBLER_PORT (or just PORT) = nibbler.port in JSON, etc
- NIBBLER_DIRECTORY_STATIC = nibbler.directory.static in JSON, etc, defaults to "/public" 
- NIBBLER_AC_ALLOW_ORIGIN = nibbler.ac.allow.origin in JSON, etc, defaults to "*" 
- NIBBLER_AC_ALLOW_HEADERS = nibbler.ac.allow.headers in JSON, etc, defaults to "Origin, Accept, Accept-Version, 
Content-Length, Content-MD5, Content-Type, Date, X-Api-Version, X-Response-Time, X-PINGOTHER, X-CSRF-Token, Authorization"
- NIBBLER_AC_ALLOW_METHODS = nibbler.ac.allow.methods in JSON, etc, defaults to "GET, POST, OPTIONS, PUT, PATCH, DELETE" 

For specific configuration values for a given extension, look at the relevant module README.md.

A sample config example is provided "./sample/config.json":

```json
{
  "port": 5001,
  "nibbler": {
    "port": 5000
  }
}
```

There are a few things to notice:

- both port and nibbler.port are defined.  The nibbler.port value takes priority over the port value (it's considered 
more specific).  So you will notice when the app starts that it uses nibbler.port's value when started.  The "port" 
value (really, the PORT env var) is frequently used by PaaS providers and is often a requirement for apps and containers 
to register as "healthy" or "started".  In any case, nibbler.port is also more explicit, so it takes priority.  To 
experiment, set the environment variable NIBBLER_PORT to something like 8001 and start the app.  You'll see that 
environment variables in the sample app take priorityover file-defined values.  If you remove all environment variables 
and values in the JSON file, you'll see the server start without listening on a port.  Your app can define sources its 
own way, so keep in mind the sample is just one demonstration of how this can be done.

- the environment variables can be directly mapped to JSON values (when provided as a source to LoadConfiguration). 
Environment variables 
are all caps, and underscores are used where dots were used in the JSON format.

## Logging

A simple logger must be passed to most Nibbler methods.  Some simple logger implementations have been provided:

- DefaultLogger - logs to the console
- SilentLogger - logs nothing 

## Build utilities

A few optional build utilities are included.  Specifically, ./build contains a simplistic way to cross-compile as well
as embed the git tag into your app.

To use, make something like build/main.go in your app, which contains:

```go
import "github.com/markdicksonjr/nibbler/build"
...
build.ProcessDefaultTargets("BinaryName", "main/main.go")
```

This will build your app for a few platforms (currently, Windows, darwin, linux).

To access the git tag in your app, use:

```go
import "github.com/markdicksonjr/nibbler/build"
...
build.GitTag
```

## Auto-wiring

To prepare for very complex apps, an auto-wiring mechanism has been added.  This mechanism will take a given slice of 
allocated extensions and assign Extension pointer field values for each of them that is undefined, as well as order the 
extensions for initialization.

There are currently a few restrictions, however.  The current auto-wiring implementation still has some trouble where 
fields are interfaces.  If the field is a pointer to a struct type (e.g. *sendgrid.Extension), the auto-wiring will 
work fine.  

Extensions like the user extension are currently trouble for auto-wiring, as the extension has a field that is an 
interface (which is also an extension).  For now, manually wire something like this, and automatically wire everything 
else.

Example:

```
extensions := []nibbler.Extension{
    &A{},
    &A1{},
    &B1{},
    &AB{},
    &B{},
    &C{},
    &BC{},
}
exts, err := nibbler.AutoWireExtensions(extensions, logger)

// check error

err = app.Init(config, logger, extensions)

// check error
```