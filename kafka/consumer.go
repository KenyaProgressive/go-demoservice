package kafka

import (
	"context"
	"go-demoservice/utils"

	"github.com/segmentio/kafka-go"
)

func consumerLoop(ctx context.Context, messageReader *kafka.Reader) error {
	for {
		_, err := messageReader.ReadMessage(ctx)
		if err != nil {
			utils.KafkaReadLogger.Error(err)
			return err
		}
	}
}
