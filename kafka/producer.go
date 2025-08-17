package mykafka

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go-demoservice/utils"
	"strconv"
	"sync"

	"github.com/ddosify/go-faker/faker"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func GenerateMessages(encoder *json.Encoder, writer *kafka.Writer, buff *bytes.Buffer, wg *sync.WaitGroup, cacheMap map[string]utils.Message, ctx context.Context) {
	utils.BaseLogger.Info("Producer launched sucessfully")

	message := utils.Message{}
	messageIndices := make([]int, 20)
	maxMessages := len(messageIndices)

	faker := faker.NewFaker()

	for i := range messageIndices {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			buff.Reset()
			makeMessage(&message, &faker)
			if err := encoder.Encode(message); err != nil {
				utils.BaseLogger.Errorf("Message %d wasn't encoded to JSON with error: %s", i, err)
				continue
			}
			err := writer.WriteMessages(ctx, kafka.Message{
				Value: buff.Bytes(),
				Key:   []byte("OrdersKey")})
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					utils.KafkaWriteLogger.Error("The time to send the message exceeded the allowed value")
				} else {
					utils.KafkaWriteLogger.Error(err)
				}
				continue
			}
			utils.KafkaWriteLogger.Infof("Message with uuid %s succesfully pushed in broker", message.OrderUId)

			if i > maxMessages-10 {
				// last 10 messages will be in cache
				cacheMap[message.OrderUId] = message
			}
		}
	}
}

func makeMessage(message *utils.Message, randInfoGen *faker.Faker) {
	paymentTimestamp, createdDateUnix := randInfoGen.CurrentTimestamp(), randInfoGen.CurrentISOTimestamp()

	message.OrderUId = uuid.NewString()
	message.TrackNumber = "WBILMTESTTRACK"
	message.Entry = "WBIL"
	message.Delivery = utils.Delivery{
		Name:    randInfoGen.RandomPersonFullName(),
		Phone:   randInfoGen.RandomPhoneNumber(),
		ZipCode: utils.TestMessageZipCodes[randInfoGen.RandomIntBetween(0, 9)],
		City:    randInfoGen.RandomAddressCity(),
		Address: randInfoGen.RandomAddressStreetAddress(),
		Region:  randInfoGen.RandomAddressCountry(),
		Email:   randInfoGen.RandomEmail()}
	message.Payment = utils.Payment{
		Transaction:  message.OrderUId,
		RequestId:    "",
		Currency:     randInfoGen.RandomCurrencyCode(),
		Provider:     "wbpay",
		Amount:       0, // Counting will be below this big struct
		PaymentDt:    uint(paymentTimestamp),
		Bank:         utils.TestMessageBankNames[randInfoGen.RandomIntBetween(0, 9)],
		DeliveryCost: uint(randInfoGen.RandomIntBetween(0, 7000)),
		GoodsTotal:   0, // Counting will be below this big struct
		CustomFee:    0}
	message.Items = []utils.Items{
		{
			ChrtID:      uint(randInfoGen.RandomIntBetween(0, 999999999)),
			TrackNumber: "WBILMTESTTRACK",
			Price:       uint(randInfoGen.RandomInt()),
			Rid:         uuid.NewString(),
			Name:        randInfoGen.RandomProductName(),
			Sale:        uint(randInfoGen.RandomIntBetween(0, 70)),
			Size:        strconv.Itoa(randInfoGen.RandomIntBetween(32, 60)),
			TotalPrice:  uint(randInfoGen.RandomIntBetween(100, 8000)),
			NmID:        uint(randInfoGen.RandomIntBetween(0, 999999999)),
			Brand:       randInfoGen.RandomWord(),
			Status:      randInfoGen.RandomInt(),
		}}
	message.Locale = randInfoGen.RandomLocale()
	message.InternalSignature = ""
	message.CustomerId = strconv.Itoa(randInfoGen.RandomInt()*randInfoGen.RandomInt()) + "test"
	message.DeliveryService = "meest"
	message.ShardKey = strconv.Itoa(randInfoGen.RandomIntBetween(1, 10))
	message.SmId = uint(randInfoGen.RandomInt())
	message.DateCreated = createdDateUnix
	message.OofShard = strconv.Itoa(randInfoGen.RandomIntBetween(1, 10))

	goodsTotalValue := 0

	for _, item := range message.Items {
		goodsTotalValue += int(item.TotalPrice)
	}
	message.Payment.GoodsTotal = uint(goodsTotalValue)
	message.Payment.Amount = message.Payment.GoodsTotal + message.Payment.DeliveryCost
}

// func PhoneNumberValidation(num string) string {
// 	noNumber := "NO_PHONE_NUMBER"
// 	pn, err := phonenumbers.Parse(num, "")
// 	if err != nil {
// 		return noNumber
// 	}

// 	if phonenumbers.IsValidNumber(pn) {
// 		return pn.String()
// 	}

// 	return noNumber
// }
