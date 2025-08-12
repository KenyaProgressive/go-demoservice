package main

import (
	"go-demoservice/db"
	"go-demoservice/utils"

	_ "github.com/segmentio/kafka-go"
)

func main() {
	db, err := db.MakeDbConnection()
	// messageReader := kafka.NewReader(utils.KafkaReaderConfig)
	// defer messageReader.Close()

	if err != nil {
		utils.DbLogger.Error(err)
		panic("")
	}
	defer db.Close()
}
