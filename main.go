package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-demoservice/db"
	mykafka "go-demoservice/kafka"
	"go-demoservice/utils"
	"sync"
	_ "sync"

	"github.com/segmentio/kafka-go"
)

func main() {
	dbase, err := db.MakeDbConnection()
	if err != nil {
		utils.DbLogger.Error(err)
		panic("")
	}
	defer dbase.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	utils.BaseLogger.Info("Start working of a service")

	messageWriter := kafka.NewWriter(utils.KafkaWriterConfig)
	messageReader := kafka.NewReader(utils.KafkaReaderConfig)

	defer messageReader.Close()
	defer messageWriter.Close()

	buff := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buff)

	go mykafka.GenerateMessages(encoder, messageWriter, buff, &wg)
	go mykafka.ConsumerLoop(messageReader, &wg, dbase)

	wg.Wait()

	fmt.Println("Test service successfully completed work")
	utils.BaseLogger.Info("Service successfully exited without errors")
}
