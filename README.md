# Nibbler

[![CircleCI](https://circleci.com/gh/markdicksonjr/nibbler.svg?style=svg)](https://circleci.com/gh/markdicksonjr/nibbler)

An extensible group of modules designed to take a lot of the boilerplate out of the code for a top-notch
Go web server.  These modules require Go v1.9+.  

At first glance, Nibbler's code-base seems to be quite monolithic, but by leveraging the capabilities of Go's import 
mechanisms, your web app can easily get and import only what is needed.

## Module Categories

Modules in Nibbler are organized by functional category, with the main Nibbler structures at the root level.  These are the 
module categories below the root level:

- Auth - authentication/authorization modules that do not integrate with Nibbler's user model (source of truth is not Nibbler).
- Database - connect to databases and expose mechanisms to create, query, etc.
- Mail - both inbound and outbound email/sms/etc
- Session - session storage and retreival
- Storage - block/blob storage
- User - the Nibbler user model, and various integrations that can operate with it.  These will tend to be auth integrations.

## Running the included sample apps

Sample apps have been provided to show how Nibbler and some extensions are used.  They'll be helpful as I fill in documentation gaps.

First, grab dependencies.  If using dep, it will be with something like

`dep ensure`

or you could use something like

`go get`

then, build the sample app (from the correct directory):

`go build`

finally, run the app (from ./sample) with

`go run sample.application`

Note that using dep this way could pull in more vendor dependencies than your app might need (e.g. elasticsearch-related
dependencies when you're using only SQL).  To avoid this, use go get (perhaps with the Gopkg.toml as a reference for the full
dependency list).

## Configuration

For specific configuration values for a given extension, look at the relevant module README.md.

The app configuration can be created with this method:

```go
config, err := nibbler.LoadApplicationConfiguration(nil)
```

If nil is provided to the core.LoadApplicationConfiguration method, it will use environment variables for
configuration of your app.  This feature can be overridden by doing something like this in your app:

```go
envSources := []source.Source{
    file.NewSource(file.WithPath("./sample.config.json")),
    env.NewSource(),
}

config, err := nibbler.LoadApplicationConfiguration(&envSources)
```

In this case, it will first apply environment variables, then apply properties from the json file (reverse order of definition). 
The environment variables used are upper-case but with dot-notation where the dots are replaced by underscores (e.g. nibbler.port is NIBBLER_PORT).

The following properties are available by default, but custom properties can be obtained from your extension from the Application
password to it in the Init() method.  For example:

```
(*app.GetConfiguration().Raw).Get("some", "property")
```

The core properties (all optional) are:

- NIBBLER_PORT (or just PORT) = nibbler.port in JSON, etc
- NIBBLER_DIRECTORY_STATIC = nibbler.directory.static in JSON, etc
- NIBBLER_AC_ALLOW_ORIGIN = nibbler.ac.allow.origin in JSON, etc
- NIBBLER_AC_ALLOW_HEADERS = nibbler.ac.allow.headers in JSON, etc
- NIBBLER_AC_ALLOW_METHODS = nibbler.ac.allow.methods in JSON, etc

By default, nibbler will only pay attention to environment variables, but the sample application that ships with
Nibbler shows how one might apply both environment variables and files.

A sample config example is provided "sample.config.json".  There are a few things to notice:

- both port and nibbler.port are defined.  The nibbler.port value takes priority over the port value.  So you will notice when the app starts
that it uses nibbler.port's value when started.  The "port" value (really, the PORT env var) is frequently used by PaaS providers and is often a requirement for apps and 
containers to register as "healthy" or "started".  In any case, nibbler.port is also more explicit, so it takes priority.  To experiment, set the
environment variable NIBBLER_PORT to something like 8001 and start the app.  You'll see that environment variables in the sample app take priority
over file-defined values.  If you remove all environment variables and values in the JSON file, you'll see the server start on port 3000.
Your app can define sources its own way, so keep in mind the sample is just one demonstration of how this can be done.

- the environment variables can be directly mapped to JSON values (when provided as a source to LoadApplicationConfiguration).  Environment variables 
are all caps, and underscores are used where dots were used in the JSON format.