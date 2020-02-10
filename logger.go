package nibbler

import (
	"log"
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
