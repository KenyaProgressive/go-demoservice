package mykafka

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	mydb "go-demoservice/db"
	"go-demoservice/utils"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

func ConsumerLoop(messageReader *kafka.Reader, wg *sync.WaitGroup, db *sql.DB) {

	loopCtx, cancelLoop := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancelLoop()

	for {
		select {
		case <-loopCtx.Done():
			utils.KafkaReadLogger.Info("Consumer was stoped by ending TimeContext")
			wg.Done()
			return
		default:
		}

		readCtx, cancelReadCtx := context.WithTimeout(context.Background(), time.Second*30)

		msg, err := messageReader.ReadMessage(readCtx)
		cancelReadCtx()

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				utils.KafkaReadLogger.Error("The time to read the message exceeded the allowed value")
			} else {
				utils.KafkaReadLogger.Error(err)
			}
			continue // Read this message again
		}

		utils.KafkaReadLogger.Infof("Message (%s) was successfully read by consumer", string(msg.Value))

		if err := mydb.PrepareMessagesAndPushToDb(db, msg.Value); err != nil {
			utils.DbLogger.Errorf("Pushing message to database was stopped with error: %s", err)
			continue // Skipping message
		}

		fmt.Printf("-->Message with offset %d was successfully pushed\n", msg.Offset)
	}
}
