package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"go-demoservice/db"
	mykafka "go-demoservice/kafka"
	"go-demoservice/utils"
	"go-demoservice/web/backend"
	"os/signal"
	"sync"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func main() {
	dbase, err := db.MakeDbConnection()
	if err != nil {
		utils.DbLogger.Error(err)
		panic(err)
	}

	defer func() {
		if errDbConnClose := dbase.Close(); errDbConnClose != nil {
			utils.BaseLogger.Errorf("Error in Db-conn close: %s", errDbConnClose)
		}
	}()

	var wg sync.WaitGroup

	genMessageFlag := parseGenMessageFlag()

	utils.BaseLogger.Info("Start working of a service")

	messageWriter := kafka.NewWriter(utils.KafkaWriterConfig)
	messageReader := kafka.NewReader(utils.KafkaReaderConfig)

	defer func() {
		if errMsgReaderClose := messageReader.Close(); errMsgReaderClose != nil {
			utils.BaseLogger.Errorf("Error in close messaageReader: %s", errMsgReaderClose)
		}
		if errMsgWriterClose := messageWriter.Close(); errMsgWriterClose != nil {
			utils.BaseLogger.Errorf("Error in close messaageWriter: %s", errMsgWriterClose)
		}
	}()

	cacheMap := make(map[string]utils.Message)

	buff := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buff)

	globalContext, globContextCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer globContextCancel()

	if genMessageFlag {
		wg.Add(1)
		go func() {
			mykafka.GenerateMessages(encoder, messageWriter, buff, &wg, cacheMap, globalContext)
		}()
	} else {
		utils.BaseLogger.Info("Producer wasn't launch (-gen=False)")
	}

	wg.Add(1)
	go func() {
		mykafka.ConsumerLoop(messageReader, &wg, dbase, globalContext)
	}()

	wg.Add(1)
	go func() {
		backend.App(dbase, cacheMap, &wg, globalContext)
	}()

	wg.Wait()

	fmt.Println("Test service successfully completed work")
	utils.BaseLogger.Info("Service successfully exited without errors")
}

func parseGenMessageFlag() bool {
	genMessageFlag := flag.Bool("gen", false, "Generate messages via launching producer-goroutine")
	flag.Parse()

	return *genMessageFlag
}
