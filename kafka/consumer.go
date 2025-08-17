package mykafka

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	mydb "go-demoservice/db"
	"go-demoservice/utils"
	"sync"

	"github.com/segmentio/kafka-go"
)

func ConsumerLoop(messageReader *kafka.Reader, wg *sync.WaitGroup, db *sql.DB, ctx context.Context) {

	utils.BaseLogger.Info("Consumer launched successfully")

	for {
		select {
		case <-ctx.Done():
			utils.KafkaReadLogger.Info("Consumer was stoped by ending TimeContext")
			wg.Done()
			return
		default:
			msg, err := messageReader.ReadMessage(ctx)

			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					utils.KafkaReadLogger.Error("The time to read the message exceeded the allowed value")
				} else {
					utils.KafkaReadLogger.Error(err)
				}
				continue // Read this message again
			}

			if err := mydb.PrepareMessagesAndPushToDb(db, msg.Value); err != nil {
				utils.DbLogger.Errorf("Pushing message to database was stopped with error: %s", err)
				continue // Skipping message
			}

			fmt.Printf("-->Message with offset %d was successfully pushed\n", msg.Offset)
		}
	}
}
