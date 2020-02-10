package nibbler

import (
	"log"
	"os"
	"strings"
)

// DefaultLogger logs using Go's standard "log" package, always via Println
type DefaultLogger struct{}

func (logger DefaultLogger) Error(message ...string) {
	log.Println(strings.Join(message, ""))
}

func (logger DefaultLogger) Warn(message ...string) {
	log.Println(strings.Join(message, ""))
}

func (logger DefaultLogger) Info(message ...string) {
	log.Println(strings.Join(message, ""))
}

func (logger DefaultLogger) Trace(message ...string) {
	log.Println(strings.Join(message, ""))
}

func (logger DefaultLogger) Debug(message ...string) {
	log.Println(strings.Join(message, ""))
}

// SilentLogger simply doesn't do anything when logs are to be written
type SilentLogger struct{}

func (logger SilentLogger) Error(message ...string) {
}

func (logger SilentLogger) Warn(message ...string) {
}

func (logger SilentLogger) Info(message ...string) {
}

func (logger SilentLogger) Debug(message ...string) {
}

func (logger SilentLogger) Trace(message ...string) {
}

// LogFatalNonNil is a utility function that handles a conditional check against an error in such a way that it treats
// non-nil errors as fatal.  Specifically, it exits with a non-zero code
func LogFatalNonNil(logger Logger, err error, wrap ...string) {
	if err != nil {
		if len(wrap) > 0 {
			logger.Error(wrap[0] + ", " + err.Error())
		} else {
			logger.Error(err.Error())
		}
		os.Exit(1)
	}
}

// LogErrorNonNil is a utility function that allows for a conditional check against an error in such a way that it logs
// the error if it is non-nil.  It also allows for an optional "wrap" string to provide extra context to the message
func LogErrorNonNil(logger Logger, err error, wrap ...string) error {
	if err != nil {
		if len(wrap) > 0 {
			logger.Error(wrap[0] + ", " + err.Error())
		} else {
			logger.Error(err.Error())
		}
	}
	return err
}
