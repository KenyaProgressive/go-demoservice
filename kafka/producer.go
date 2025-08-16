package mykafka

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go-demoservice/utils"
	"strconv"
	"sync"
	"time"

	"github.com/ddosify/go-faker/faker"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func GenerateMessages(encoder *json.Encoder, writer *kafka.Writer, buff *bytes.Buffer, wg *sync.WaitGroup) {
	message := utils.Message{}
	messageIndices := make([]int, 10)

	faker := faker.NewFaker()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for i := range messageIndices {
		buff.Reset()
		makeMessage(&message, i, &faker)
		if err := encoder.Encode(message); err != nil {
			utils.BaseLogger.Errorf("Message %d wasn't encoded to JSON with error: %s", i, err)
			continue
		}
		err := writer.WriteMessages(ctx, kafka.Message{
			Value: buff.Bytes(),
			Key:   []byte("mykey")})
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				utils.KafkaWriteLogger.Error("The time to send the message exceeded the allowed value")
			} else {
				utils.KafkaWriteLogger.Error(err)
			}
			continue
		}
		utils.KafkaWriteLogger.Info("Message succesfully pulled in broker")
		utils.DbLogger.Debugf("PUSHED: %s", buff.String())

	}
	wg.Done()
}

func makeMessage(message *utils.Message, offset int, randInfoGen *faker.Faker) {
	paymentTimestamp, createdDateUnix := randInfoGen.CurrentTimestamp(), randInfoGen.CurrentISOTimestamp()

	message.OrderUId = uuid.NewString()
	message.TrackNumber = "WBILMTESTTRACK"
	message.Entry = "WBIL"
	message.Delivery = utils.Delivery{
		Name:    randInfoGen.RandomPersonFullName(),
		Phone:   randInfoGen.RandomPhoneNumber(),
		ZipCode: utils.TestMessageZipCodes[offset],
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
		Bank:         utils.TestMessageBankNames[offset],
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
