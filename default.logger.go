package nibbler

import "log"

type DefaultLogger struct {}

func (logger DefaultLogger) Error(message string) {
	log.Println(message)
}

func (logger DefaultLogger) Warn(message string) {
	log.Println(message)
}

func (logger DefaultLogger) Info(message string) {
	log.Println(message)
}

func (logger DefaultLogger) Debug(message string) {
	log.Println(message)
}
