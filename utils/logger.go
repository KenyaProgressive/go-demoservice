package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var BaseLogger = initBaseLogger()
var DbLogger = initDbLogger()
var KafkaReadLogger = initKafkaReadLogger()
var KafkaWriteLogger = initKafkaWriteLogger()

func initBaseLogger() *logrus.Logger {
	bl := logrus.New()
	bl.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: logTimeFormat,
		PrettyPrint:     true})

	bl.SetLevel(logrus.DebugLevel)

	logfile, err := os.OpenFile("logs/baselog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		bl.Error(err)
	} else {
		bl.SetOutput(logfile)
	}
	return bl
}

func initDbLogger() *logrus.Logger {
	lDb := logrus.New()
	lDb.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: logTimeFormat,
		PrettyPrint:     true})

	lDb.SetLevel(logrus.DebugLevel)
	logDbFile, err := os.OpenFile("logs/dblog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		lDb.Error(err)
	} else {
		lDb.SetOutput(logDbFile)
	}

	return lDb
}

func initKafkaReadLogger() *logrus.Logger {
	kfRL := logrus.New()
	kfRL.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: logTimeFormat,
		PrettyPrint:     true})

	kfRL.SetLevel(logrus.DebugLevel)

	logKfFile, err := os.OpenFile("logs/kafka-readlog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		kfRL.Error(err)
	} else {
		kfRL.SetOutput(logKfFile)
	}

	return kfRL

}

func initKafkaWriteLogger() *logrus.Logger {
	kfWL := logrus.New()
	kfWL.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: logTimeFormat,
		PrettyPrint:     true})

	kfWL.SetLevel(logrus.DebugLevel)

	logKfFile, err := os.OpenFile("logs/kafka-writelog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		kfWL.Error(err)
	} else {
		kfWL.SetOutput(logKfFile)
	}

	return kfWL

}
