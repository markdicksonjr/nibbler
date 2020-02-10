package nibbler

import (
	"log"
	"os"
	"strings"
)

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
