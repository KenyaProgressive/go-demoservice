package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

func setEnvVars() (map[string]string, error) {
	// Setting enviroment variables to global map

	if err := godotenv.Load(".env"); err != nil {
		BaseLogger.Error("NO LOADED .ENV FILE -- PLEASE DO IT BEFORE LAUNCH")
		return nil, err
	}
	envMap := map[string]string{
		"dbUserName": os.Getenv("DB_USERNAME"),
		"dbPassword": os.Getenv("DB_PASSWORD"),
		"dbHost":     os.Getenv("DB_HOST"),
		"dbPort":     os.Getenv("DB_PORT"),
		"dbName":     os.Getenv("DATABASE_NAME")}

	return envMap, nil
}

var envMap, _ = setEnvVars()
var ConnectString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", envMap["dbUserName"], envMap["dbPassword"], envMap["dbHost"], envMap["dbPort"], envMap["dbName"])
var topicName = "demoservice-orders"
var kafkaGroupID = "demoservice-group"

var KafkaWriterConfig = kafka.WriterConfig{
	Brokers:      []string{"localhost:9092"},
	Topic:        topicName,
	Balancer:     &kafka.Hash{}, // Отправка по партициям согласно хэшу указанного ключа
	RequiredAcks: 1,             // Подтверждение только от leader-partition
	Logger:       KafkaWriteLogger}

var KafkaReaderConfig = kafka.ReaderConfig{
	Brokers:  []string{"localhost:9092"},
	Topic:    topicName,
	GroupID:  kafkaGroupID,
	MaxBytes: 10e6,
	Logger:   KafkaReadLogger}
